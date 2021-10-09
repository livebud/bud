package npm

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
)

func Install(root string, packages ...string) error {
	eg := new(errgroup.Group)
	for _, pkg := range packages {
		pkg := pkg
		eg.Go(func() error {
			return install(root, pkg)
		})
	}
	return eg.Wait()
}

func resolvePackage(pkgname string) (resolved *npmPackage, err error) {
	parts := strings.Split(pkgname, "@")
	name := parts[0]
	version := "latest"
	if len(parts) >= 2 {
		version = parts[1]
	}
	if version == "latest" {
		return nil, fmt.Errorf("npm %s: latest is not supported yet", pkgname)
	}
	return &npmPackage{
		Name:    name,
		Version: version,
	}, nil
}

type npmPackage struct {
	Name    string
	Version string
}

func (p *npmPackage) URL() string {
	return fmt.Sprintf(`https://registry.npmjs.org/%[1]s/-/%[1]s-%[2]s.tgz`, p.Name, p.Version)
}

func (p *npmPackage) Dir(root string) string {
	return filepath.Join(root, "node_modules", p.Name)
}

func install(root string, pkgname string) error {
	pkg, err := resolvePackage(pkgname)
	if err != nil {
		return err
	}
	res, err := http.Get(pkg.URL())
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("npm install %s: unexpected status code %d", pkg.Name, res.StatusCode)
	}
	gzipReader, err := gzip.NewReader(res.Body)
	if err != nil {
		return err
	}
	defer gzipReader.Close()
	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		fileInfo := header.FileInfo()
		dir := filepath.Join(pkg.Dir(root), rootless(filepath.Dir(header.Name)))
		filename := filepath.Join(dir, fileInfo.Name())
		if err = os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		if _, err := io.Copy(file, tarReader); err != nil {
			return err
		}
		if err = file.Close(); err != nil {
			return err
		}
	}
	return nil
}

func rootless(fpath string) string {
	parts := strings.Split(fpath, string(filepath.Separator))
	return path.Join(parts[1:]...)
}

func Link(from string, to string) error {
	npm, err := exec.LookPath("npm")
	if err != nil {
		return err
	}
	absFrom, err := filepath.Abs(from)
	if err != nil {
		return err
	}
	absTo, err := filepath.Abs(to)
	if err != nil {
		return err
	}
	cmd := exec.Command(npm, "link", absFrom)
	cmd.Dir = absTo
	cmd.Env = os.Environ()
	stderr := new(bytes.Buffer)
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("npm link %q: %w\n\n%s", from, err, stderr)
	}
	return nil
}
