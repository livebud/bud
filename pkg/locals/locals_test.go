package locals_test

import (
	"context"
	"sync"
	"testing"

	"github.com/livebud/bud/pkg/locals"
	"github.com/matryer/is"
)

func TestSetAndGet(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	key := "testKey"
	value := "testValue"

	// Set value
	ctx = locals.Set(ctx, key, value)

	// Get value
	actual, ok := locals.Get[string](ctx, key)
	is.True(ok)             // should be true
	is.Equal(actual, value) // values should be equal
}

func TestGetNonExistentKey(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	_, ok := locals.Get[string](ctx, "nonexistent")
	is.True(!ok) // should be false
}

func TestTypeMismatchOnGet(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	key := "testKey"
	locals.Set(ctx, key, "stringValue")

	_, ok := locals.Get[int](ctx, key)
	is.True(!ok) // should be false, type mismatch
}

func TestConcurrentAccess(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	key := "concurrentKey"
	wg := sync.WaitGroup{}

	// Write in goroutines
	for i := 0; i < 100; i++ {
		wg.Add(1)
		ctx := ctx
		go func(i int) {
			defer wg.Done()
			ctx = locals.Set(ctx, key, i)
			value, ok := locals.Get[int](ctx, key)
			is.True(ok)                        // should be true
			is.True(value >= 0 && value < 100) // value should be within range
		}(i)
	}
	wg.Wait()
}

func TestOverwriteExistingKey(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	key := "overwriteKey"
	firstValue := "firstValue"
	secondValue := "secondValue"

	ctx = locals.Set(ctx, key, firstValue)
	ctx = locals.Set(ctx, key, secondValue)

	actual, ok := locals.Get[string](ctx, key)
	is.True(ok)                   // should be true
	is.Equal(actual, secondValue) // should be the second value
}
