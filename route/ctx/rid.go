package ctx

import (
	"net/http"

	helpers "github.com/robitx/inceptus/helpers"
)

type ctxKeyRequestID struct{}

// GetRequestID reads requestID from context, header or generates new one
func GetRequestID(headerName string, r *http.Request) string {
	// check context
	if requestID, ok := getStr(ctxKeyRequestID{}, r); ok {
		return requestID
	}
	// check header
	if requestID := r.Header.Get(headerName); requestID != "" {
		return requestID
	}
	// generate new ID if necessery
	return helpers.GenerateID()
}

// SetRequestID stores requestID in context
func SetRequestID(id string, r *http.Request) *http.Request {
	return set(ctxKeyRequestID{}, id, r)
}
