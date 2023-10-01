package view

import "context"

type Context map[string]any

var contextKey = struct{}{}

// SetContext sets the fields in the context
func SetContext(ctx context.Context, fields Context) context.Context {
	contextMap, ok := ctx.Value(contextKey).(Context)
	if !ok {
		contextMap = Context{}
	}
	for k, v := range fields {
		contextMap[k] = v
	}
	return context.WithValue(ctx, contextKey, contextMap)
}

// GetContext gets fields from the context
func GetContext(ctx context.Context) Context {
	contextMap, ok := ctx.Value(contextKey).(Context)
	if !ok {
		return Context{}
	}
	return contextMap
}
