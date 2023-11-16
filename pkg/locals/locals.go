package locals

import (
	"context"
)

type contextKey struct{}

var mapKey = contextKey{}

type anyMap = map[string]any

// Set a local value. Setting values is concurrency-safe.
func Set(ctx context.Context, key string, value any) context.Context {
	oldMap, ok := ctx.Value(mapKey).(anyMap)
	if !ok {
		return context.WithValue(ctx, mapKey, anyMap{
			key: value,
		})
	}
	// Make a copy for the new context
	newMap := make(anyMap, len(oldMap)+1)
	for k, v := range oldMap {
		newMap[k] = v
	}
	newMap[key] = value
	return context.WithValue(ctx, mapKey, newMap)
}

// Get a local value. Getting values is concurrency-safe.
func Get[Value any](ctx context.Context, key string) (v Value, ok bool) {
	m, ok := ctx.Value(mapKey).(anyMap)
	if !ok {
		return v, false
	}
	val, ok := m[key]
	if !ok {
		return v, false
	}
	v, ok = val.(Value)
	if !ok {
		return v, false
	}
	return v, true
}
