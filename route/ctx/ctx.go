package ctx

import (
	"context"
	"net/http"
)

// generic context setter
func set(key interface{}, value interface{}, r *http.Request) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, key, value)
	return r.WithContext(ctx)
}

// string context getter
func getStr(key interface{}, r *http.Request) (string, bool) {
	ctx := r.Context()
	val, ok := ctx.Value(key).(string)
	return val, ok
}
