package testdir

import (
	"io/fs"
	"strings"
	"testing/fstest"

	"github.com/lithammer/dedent"
	"gitlab.com/mnm/bud/pkg/modcache"
)

func redent(s string) string {
	return strings.TrimSpace(dedent.Dedent(s)) + "\n"
}

var goMod = redent(`
	module app.com

	require (
		gitlab.com/mnm/bud v0.0.0
	)
`)

// func hash(input string) string {
// 	hash := xxhash.New()
// 	hash.Write([]byte(input))
// 	return base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
// }

func New() *Dir {
	return &Dir{
		Modules:     map[string]modcache.Files{},
		NodeModules: map[string]string{},
		BFiles:      map[string][]byte{},
		Files:       map[string]string{},
	}
}

type Dir struct {
	Files       map[string]string         // String files (convenient)
	BFiles      map[string][]byte         // Byte files (for images and binaries)
	Modules     map[string]modcache.Files // name@version[path[data]]
	NodeModules map[string]string         // name[version]
}

func (d *Dir) fileSystem() (fs.FS, error) {
	mapfs := fstest.MapFS{}
	// for pv, files := range d.Modules {
	// parts := strings.Split(pv, "@")
	// if len(parts) != 2 {
	// 	return nil, fmt.Errorf("Modules must have path@version key")
	// }
	// path, version := parts[0], parts[1]
	// for _, file := range files {

	// }
	// }
	// if len(d.Modules) > 0 {

	// 	mapfs[]
	// // if len(d.Modules) > 0 {
	// // 	cacheDir := g.t.TempDir()
	// // 	modCache = modcache.New(cacheDir)
	// // 	// Importing from the snapshot
	// // 	if _, err := os.Stat(modSnapDir); nil == err {
	// // 		if err := modCache.Import(modSnapDir); err != nil {
	// // 			return err
	// // 		}
	// // 	}

	// }
	return mapfs, nil
}

// func (d *Dir) HashKey() (string, error) {
// 	buf, err := json.Marshal(d)
// 	if err != nil {
// 		return "", err
// 	}
// 	return hash(string(buf)), err
// }

// CacheDir returns the module cache dir based on dir
// func CacheDir(dir string) string {
// 	return filepath.Join(dir, ".mod")
// }

// // SnapshotDir returns the snapshot directory
// func SnapshotDir(dir string) (string, error) {
// 	cacheDir, err := os.UserCacheDir()
// 	if err != nil {
// 		return "", err
// 	}
// 	return filepath.Join(cacheDir, "bud", "testdir"), nil
// }

func (d *Dir) Write(dir string) error {
	fsys, err := d.fileSystem()
	if err != nil {
		return err
	}
	_ = fsys
	// modCache := modcache.Default()
	// // Try loading the snapshot dir
	// modSnapDir, err := getModSnapDir(g.t.Name())
	// if err != nil {
	// 	return err
	// }
	// // Add modules
	// if len(d.Modules) > 0 {
	// 	cacheDir := g.t.TempDir()
	// 	modCache = modcache.New(cacheDir)
	// 	// Importing from the snapshot
	// 	if _, err := os.Stat(modSnapDir); nil == err {
	// 		if err := modCache.Import(modSnapDir); err != nil {
	// 			return err
	// 		}
	// 	}
	// 	// Generate a custom go.mod for the module
	// 	for pathVersion, module := range d.Modules {
	// 		if _, ok := module["go.mod"]; !ok {
	// 			gomod, err := moduleGoMod(pathVersion, module)
	// 			if err != nil {
	// 				return err
	// 			}
	// 			module["go.mod"] = gomod
	// 		}
	// 	}
	// 	if err := modCache.Write(g.Modules); err != nil {
	// 		return err
	// 	}
	// }
	// // Replace bud in Go mod if present
	// if code, ok := g.Files["go.mod"]; ok {
	// 	code, err := replaceBud(string(code))
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if len(g.Modules) > 0 {
	// 		code, err = addModules(code, g.Modules)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// 	g.Files["go.mod"] = []byte(code)
	// }
	// // Setup the application dir
	// appDir := filepath.Join("_tmp", g.t.Name())
	// if err := os.RemoveAll(appDir); err != nil {
	// 	return err
	// }
	// g.t.Cleanup(cleanup(g.t, "_tmp", appDir))
	// // Add node_modules
	// var nodeModules []string
	// if len(d.NodeModules) > 0 {
	// 	packageJSON := &npm.Package{
	// 		Name:         filepath.Base(appDir),
	// 		Version:      "0.0.0",
	// 		Dependencies: map[string]string{},
	// 	}
	// 	for name, version := range d.NodeModules {
	// 		nodeModules = append(nodeModules, name+"@"+version)
	// 		packageJSON.Dependencies[name] = version
	// 	}
	// 	pkgJSON, err := json.MarshalIndent(packageJSON, "", "  ")
	// 	if err != nil {
	// 		return err
	// 	}
	// 	d.Files["package.json"] = append(pkgJSON, '\n')
	// }
	// // Write the files to the application directory
	// err = vfs.Write(appDir, vfs.Map(d.Files))
	// if err != nil {
	// 	return err
	// }
	// // Install node_modules
	// if len(nodeModules) > 0 {
	// 	if err := npm.Install(appDir, nodeModules...); err != nil {
	// 		return err
	// 	}
	// }
	return nil
}
