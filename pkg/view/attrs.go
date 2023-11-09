package view

import "context"

const attrKey contextKey = "view:attr"

// Set an attribute on the context
func Set(parent context.Context, key string, value any) context.Context {
	values, ok := parent.Value(attrKey).(map[string]any)
	if !ok {
		values = map[string]any{}
	}
	values[key] = value
	return context.WithValue(parent, attrKey, values)
}

// GetAttrs gets the attributes from the context
func GetAttrs(ctx context.Context) (attrs map[string]any) {
	values, ok := ctx.Value(attrKey).(map[string]any)
	if !ok {
		return map[string]any{}
	}
	return values
}
