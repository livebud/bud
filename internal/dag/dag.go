package dag

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"sort"
	"strings"

	"github.com/livebud/bud/package/genfs"
	"github.com/livebud/bud/package/virt"
	_ "github.com/mattn/go-sqlite3"
)

var ErrNotFound = fmt.Errorf("dag: file not found")

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

func Load(fsys fs.FS, path string) (*Cache, error) {
	// Open the SQLite database
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	// Create the files and links tables.
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS files (
			path TEXT PRIMARY KEY,
			data BLOB,
			mode INTEGER
		)
	`); err != nil {
		return nil, err
	}
	if _, err := db.Exec(`
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

	return &Cache{fsys, db}, nil
}

type Cache struct {
	fsys fs.FS
	db   *sql.DB
}

var _ genfs.Cache = (*Cache)(nil)

func (c *Cache) Set(path string, file *virt.File) error {
	file.Path = path
	// Insert the file into the files table.
	if _, err := c.db.Exec(`
		INSERT INTO files (path, data, mode)
		VALUES (?, ?, ?)
	`, file.Path, file.Data, file.Mode); err != nil {
		return err
	}
	return nil
}

func (c *Cache) Get(path string) (*virt.File, error) {
	row := c.db.QueryRow(`
		SELECT path, data, mode
		FROM files
		WHERE path = ?
	`, path)
	file := new(virt.File)
	if err := row.Scan(&file.Path, &file.Data, &file.Mode); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w %q", ErrNotFound, path)
		}
		return nil, err
	}
	return file, nil
}

func (c *Cache) Link(from string, toPatterns ...string) error {
	params := make([]interface{}, len(toPatterns)*2)
	for i, to := range toPatterns {
		params[i*2] = from
		params[i*2+1] = to
	}
	// Insert the links of the file into the links table.
	sql := new(strings.Builder)
	sql.WriteString(`INSERT INTO links (from_path, to_path) VALUES `)
	for i := range toPatterns {
		if i > 0 {
			sql.WriteString(", ")
		}
		sql.WriteString(`(?, ?)`)
	}
	if _, err := c.db.Exec(sql.String(), params...); err != nil {
		return err
	}
	return nil
}

// Ancestors returns all the ancestor paths. Those that depend on the given
// paths.
func (c *Cache) Ancestors(paths ...string) (deps []string, err error) {
	visited := make(map[string]bool)
	for i := 0; i < len(paths); i++ {
		const sql = `
			SELECT from_path
			FROM links
			WHERE to_path = ?
		`
		rows, err := c.db.Query(sql, paths[i])
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var from string
			if err := rows.Scan(&from); err != nil {
				return nil, err
			}
			if !visited[from] {
				paths = append(paths, from)
				deps = append(deps, from)
				visited[from] = true
			}
		}
		if err := rows.Err(); err != nil {
			return nil, err
		} else if err := rows.Close(); err != nil {
			return nil, err
		}
	}
	sort.Strings(deps)
	return deps, nil
}

func (c *Cache) Delete(paths ...string) error {
	// Get all the ancestors paths
	ancestors, err := c.Ancestors(paths...)
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
	_, err = c.db.Exec(sql, params...)
	if err != nil {
		return err
	}

	// Delete all links to these paths
	sql = `
		DELETE FROM links
		WHERE to_path IN ` + parameterize(len(params)) + `
		OR from_path IN ` + parameterize(len(params)) + `
	`
	_, err = c.db.Exec(sql, append(params, params...)...)
	if err != nil {
		return err
	}

	return nil
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

func (c *Cache) Files() ([]*File, error) {
	rows, err := c.db.Query(`
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

func (c *Cache) Links() ([]*Link, error) {
	rows, err := c.db.Query(`
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
func (c *Cache) Print(w io.Writer) error {
	// Open the digraph
	fmt.Fprintf(w, "digraph dag {\n")
	// Print the nodes of the graph.
	rows, err := c.db.Query(`
		SELECT path
		FROM files
	`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			return err
		}
		fmt.Fprintf(w, "\t%[1]q\n", path)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	// Print the edges of the graph.
	rows, err = c.db.Query(`
		SELECT from_path, to_path
		FROM links
	`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var from, to string
		if err := rows.Scan(&from, &to); err != nil {
			return err
		}
		fmt.Fprintf(w, "\t%q -> %q\n", from, to)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	// Close the graph.
	fmt.Fprintf(w, "}\n")
	return nil
}

func (c *Cache) Close() error {
	return c.db.Close()
}
