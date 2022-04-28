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

func resolvePackage(pkgname string) (resolved *Package, err error) {
	parts := strings.Split(pkgname, "@")
	name := parts[0]
	version := "latest"
	if len(parts) >= 2 {
		version = parts[1]
	}
	if version == "latest" {
		return nil, fmt.Errorf("npm %s: latest is not supported yet", pkgname)
	}
	return newPackage(name, version), nil
}

func newPackage(name, version string) *Package {
	return &Package{
		Name:         name,
		Version:      version,
		Dependencies: map[string]string{},
	}
}

type Package struct {
	Name            string            `json:"name,omitempty"`
	Version         string            `json:"version,omitempty"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func (p *Package) URL() string {
	return fmt.Sprintf(`https://registry.npmjs.org/%[1]s/-/%[1]s-%[2]s.tgz`, p.Name, p.Version)
}

func (p *Package) Dir(root string) string {
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

// Set a some values in a package.json
func Set(dir string, kvs map[string]string) error {
	if len(kvs) == 0 {
		return fmt.Errorf("npm set: must set one or more key values")
	}
	npm, err := exec.LookPath("npm")
	if err != nil {
		return err
	}
	args := []string{"pkg", "set"}
	for key, value := range kvs {
		args = append(args, fmt.Sprintf("%s=%s", key, value))
	}
	cmd := exec.Command(npm, args...)
	cmd.Dir = dir
	stderr := new(bytes.Buffer)
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("npm %s: %w\n\n%s", strings.Join(args, " "), err, stderr)
	}
	return nil
}
