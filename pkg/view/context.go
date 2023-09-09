package view

import "context"

type Data map[string]any

type contextKey struct{}

var dataKey contextKey

func Set(ctx context.Context, data Data) context.Context {
	dataMap, ok := ctx.Value(dataKey).(Data)
	if !ok {
		dataMap = Data{}
	}
	for k, v := range data {
		dataMap[k] = v
	}
	return context.WithValue(ctx, dataKey, dataMap)
}

// get all the data from the context
func getAll(ctx context.Context) Data {
	dataMap, ok := ctx.Value(dataKey).(Data)
	if !ok {
		return Data{}
	}
	return dataMap
}
