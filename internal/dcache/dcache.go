package dcache

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var ErrNotFound = fmt.Errorf("dcache: file not found")

type File struct {
	Path  string
	Data  []byte
	Mode  fs.FileMode
	Links []string
}

type Link struct {
	From string
	To   string
}

type Cache interface {
	Get(ctx context.Context, path string) (*File, error)
	Put(ctx context.Context, path string, file *File) error
	Delete(ctx context.Context, paths ...string) error
	Close() error
}

func Load(ctx context.Context, dbPath string) (*cache, error) {
	// Open the SQLite database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create the files and links tables.
	if _, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS files (
			path TEXT PRIMARY KEY,
			data BLOB,
			mode INTEGER
		)
	`); err != nil {
		return nil, err
	}
	if _, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS links (
			from_path TEXT,
			to_path TEXT,
			FOREIGN KEY(from_path) REFERENCES files(path),
			FOREIGN KEY(to_path) REFERENCES files(path),
			PRIMARY KEY(from_path, to_path)
		)
	`); err != nil {
		return nil, err
	}

	return &cache{db}, nil
}

type cache struct {
	db *sql.DB
}

var _ Cache = (*cache)(nil)

func (c *cache) Put(ctx context.Context, path string, file *File) error {
	file.Path = path

	// Start a transaction.
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert the file into the files table.
	if _, err := tx.Exec(`
		INSERT INTO files (path, data, mode)
		VALUES (?, ?, ?)
	`, path, file.Data, file.Mode); err != nil {
		return err
	}

	// Insert the links of the file into the links table.
	for _, to := range file.Links {
		if _, err := tx.Exec(`
			INSERT INTO links (from_path, to_path)
			VALUES (?, ?)
		`, path, to); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (c *cache) Get(ctx context.Context, path string) (*File, error) {
	row := c.db.QueryRow(`
		SELECT path, data, mode
		FROM files
		WHERE path = ?
	`, path)
	var file File
	if err := row.Scan(&file.Path, &file.Data, &file.Mode); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w %q", ErrNotFound, path)
		}
		return nil, err
	}

	// Retrieve the links of the file.
	rows, err := c.db.Query(`
		SELECT to_path
		FROM links
		WHERE from_path = ?
	`, path)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var to string
		if err := rows.Scan(&to); err != nil {
			return nil, err
		}
		file.Links = append(file.Links, to)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &file, nil
}

// Ancestors returns all the ancestor paths. Those that depend on the given
// paths.
func (c *cache) Ancestors(ctx context.Context, paths ...string) (deps []string, err error) {
	// Start a transaction.
	tx, err := c.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	deps, err = c.ancestors(ctx, tx, paths...)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return deps, nil
}

func (c *cache) ancestors(ctx context.Context, tx *sql.Tx, paths ...string) (deps []string, err error) {
	params := make([]interface{}, len(paths))
	for i, path := range paths {
		params[i] = path
	}
	// Use a recursive CTE to retrieve the deps of the file.
	sql := `
		WITH RECURSIVE deps AS (
			SELECT from_path
			FROM links
			WHERE to_path IN ` + parameterize(len(params)) + `
			UNION ALL
			SELECT links.from_path
			FROM links
			JOIN deps
			ON links.to_path = deps.from_path
		)
		SELECT DISTINCT from_path
		FROM deps
		ORDER BY from_path
	`
	rows, err := tx.QueryContext(ctx, sql, params...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var to string
		if err := rows.Scan(&to); err != nil {
			return nil, err
		}
		deps = append(deps, to)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return deps, nil
}

func (c *cache) Delete(ctx context.Context, paths ...string) error {
	// Start a transaction
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get all the ancestors paths
	ancestors, err := c.ancestors(ctx, tx, paths...)
	if err != nil {
		return err
	}
	params := make([]interface{}, len(paths)+len(ancestors))
	for i, path := range paths {
		params[i] = path
	}
	for i, path := range ancestors {
		params[i+len(paths)] = path
	}

	// Delete all paths and their dependencies
	sql := `
		DELETE FROM files
		WHERE path IN ` + parameterize(len(params)) + `
	`
	_, err = tx.ExecContext(ctx, sql, params...)
	if err != nil {
		return err
	}

	// Delete all links to these paths
	sql = `
		DELETE FROM links
		WHERE to_path IN ` + parameterize(len(params)) + `
		OR from_path IN ` + parameterize(len(params)) + `
	`
	_, err = tx.ExecContext(ctx, sql, append(params, params...)...)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func parameterize(n int) string {
	out := "("
	for i := 0; i < n; i++ {
		if i > 0 {
			out += ", "
		}
		out += "?"
	}
	out += ")"
	return out
}

func (c *cache) Files(ctx context.Context) ([]*File, error) {
	rows, err := c.db.QueryContext(ctx, `
		SELECT path, data, mode
		FROM files
		ORDER BY path
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var files []*File
	for rows.Next() {
		var file File
		if err := rows.Scan(&file.Path, &file.Data, &file.Mode); err != nil {
			return nil, err
		}
		files = append(files, &file)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return files, nil
}

func (c *cache) Links(ctx context.Context) ([]*Link, error) {
	rows, err := c.db.QueryContext(ctx, `
		SELECT from_path, to_path
		FROM links
		ORDER BY from_path, to_path
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var links []*Link
	for rows.Next() {
		var link Link
		if err := rows.Scan(&link.From, &link.To); err != nil {
			return nil, err
		}
		links = append(links, &link)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return links, nil
}

// Print the DAG in a dot graph format. That you can paste in here:
// https://dreampuf.github.io/GraphvizOnline
func (c *cache) Print(ctx context.Context) (string, error) {
	out := new(strings.Builder)
	// Open the digraph
	out.WriteString("digraph dag {\n")
	// Print the nodes of the graph.
	rows, err := c.db.Query(`
		SELECT path
		FROM files
	`)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			return "", err
		}
		fmt.Fprintf(out, "\t%[1]q\n", path)
	}
	if err := rows.Err(); err != nil {
		return "", err
	}
	// Print the edges of the graph.
	rows, err = c.db.Query(`
		SELECT from_path, to_path
		FROM links
	`)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		var from, to string
		if err := rows.Scan(&from, &to); err != nil {
			return "", err
		}
		fmt.Fprintf(out, "\t%q -> %q\n", from, to)
	}
	if err := rows.Err(); err != nil {
		return "", err
	}
	// Close the graph.
	out.WriteString("}\n")
	return out.String(), nil
}

func (c *cache) Close() error {
	return c.db.Close()
}
