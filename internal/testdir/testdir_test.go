package testdir_test

// func TestHashKey(t *testing.T) {
// 	is := is.New(t)
// 	dt := testdir.New()
// 	dt.Files["main.go"] = `package main`
// 	key, err := dt.HashKey()
// 	is.NoErr(err)
// 	is.Equal(key, `4pvk587AtME`)
// 	dt.BFiles["favicon.ico"] = []byte("some bytes")
// 	key, err = dt.HashKey()
// 	is.NoErr(err)
// 	is.Equal(key, `-wiarUhIQyY`)
// 	dt.Files["main.go"] = `package main; func main() {}`
// 	key, err = dt.HashKey()
// 	is.NoErr(err)
// 	is.Equal(key, `lz1Tg7z1Gps`)
// 	dt.NodeModules["svelte"] = `3.46.4`
// 	key, err = dt.HashKey()
// 	is.NoErr(err)
// 	is.Equal(key, `IdA0AqXUhBA`)
// 	dt.Modules["gitlab.com/mnm/bud-tailwind@v0.0.1"] = modcache.Files{
// 		"public/tailwind/preflight.css": `/* tailwind */`,
// 	}
// 	key, err = dt.HashKey()
// 	is.NoErr(err)
// 	is.Equal(key, `DKlArUdy82U`)
// }
