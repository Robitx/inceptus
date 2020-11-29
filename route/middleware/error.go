package middleware

import (
	"net/http"
)

// ErrorResponseWriter ..
type errorResponseWriter struct {
	w http.ResponseWriter
	intercepted bool
	responses map[int][]byte
	contentType string
}

// making sure my ErrorResponseWriter satisfies normal ResponseWriter interface
var _ http.ResponseWriter = &errorResponseWriter{}

func (e *errorResponseWriter) Header() http.Header {
	return e.w.Header()
}

func (e *errorResponseWriter) Write(b []byte) (int, error) {
	if !e.intercepted {
		return e.w.Write(b)
	}

	// nothing if we intercepted the WriteHeader 
	return 0, nil
}

func (e *errorResponseWriter) WriteHeader(statusCode int) {
	if statusCode == http.StatusOK {
		e.w.WriteHeader(statusCode)
		return
	}
	if response, ok := e.responses[statusCode]; ok {
		e.w.Header().Set("Content-Type", e.contentType)
		e.w.Header().Set("X-Content-Type-Options", "nosniff")
		e.w.WriteHeader(statusCode)
		e.w.Write(response)
		e.intercepted = true
		return
	}

	e.w.WriteHeader(statusCode)
}

// Error allows overriding net/http error responses with your own..
//
// You might want to use one for rest api with text/plain content and
// otherone for stuff vidible by user with text/html content type.
// 
// Here is an example usage:
//		contentType := "text/plain; charset=utf-8"
//		customErrors := make(map[int][]byte)
//		customErrors[http.StatusNotFound] = []byte("custom 404..")
// 		...
//		
//		errorMiddleware := middleware.Error(contentType, customErrors)
//		server := &http.Server{
//			Addr: ":XXXX",
//			Handler: errorMiddleware(router),
//		}
//		server.ListenAndServe()
func Error(
	contentType string,
	responses map[int][]byte,
	) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ew := &errorResponseWriter{
				w: w,
				intercepted: false,
				responses: responses,
				contentType: contentType,
			}
			next.ServeHTTP(ew, r)
		}
		return http.HandlerFunc(fn)
	}
}