```go
modCache := modcache.Default()
module, err := mod.Load(dir, mod.WithCache(modCache))
singleflight := singleflight.New()
cache1 := cachefs.New(module, singleflight)
pluginfs := pluginfs.New(cache1)
cache2 := cachefs.New(module, singleflight)



```
