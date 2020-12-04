package ctx

import (
	"net/http"
)

type ctxKeyUserID struct{}

// GetUserID reads userId from context
func GetUserID(r *http.Request) string {
	// check context
	if id, ok := getStr(ctxKeyUserID{}, r); ok {
		return id
	}
	// no id
	return ""
}

// SetUserID stores userID in context
func SetUserID(id string, r *http.Request) *http.Request {
	return set(ctxKeyUserID{}, id, r)
}
