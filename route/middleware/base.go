package middleware

import (
	// "fmt"
	"net/http"
	"runtime/debug"
	"time"

	log "github.com/robitx/inceptus/log"
	c "github.com/robitx/inceptus/route/ctx"

	chi_middleware "github.com/go-chi/chi/middleware"
)

// Base middleware which
// - loggs information about request and response
// - catches panic and logs out the debug.Stack
// - prepares requestID (from specified header or generates new one)
func Base(
  logger *log.Logger,
  requestIDheader string,
  ) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
      
      requestID := c.GetRequestID(requestIDheader, r)
      r = c.SetRequestID(requestID, r)


			ww := chi_middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				t2 := time.Now()

				// Recover and record stack traces in case of a panic
				if rec := recover(); rec != nil {
          logger.Error().
						Bytes("debug_stack", debug.Stack()).
						Interface("recover_info", rec).
            Str("id", requestID).
						Msg("log system error")
          http.Error(ww,
            http.StatusText(http.StatusInternalServerError),
            http.StatusInternalServerError)
        }

        scheme := "http"
        if r.TLS != nil {
          scheme = "https"
        }
        userID := c.GetUserID(r)

				// log end request
        logger.Info().
          Float64("latency", float64(t2.Sub(t1).Nanoseconds()) / 1e9).
          Int("bytesOut", ww.BytesWritten()).
          Int("status", ww.Status()).
          Str("bytesIn", r.Header.Get("Content-Length")).
          Str("host", r.Host).
          Str("method", r.Method).
          Str("protocol", r.Proto).
          Str("remoteIP", r.RemoteAddr).
          Str("id", requestID).
          Str("scheme", scheme).
          Str("url", r.URL.Path).
          Str("uid", userID).
          Str("userAgent", r.Header.Get("User-Agent")).
					Msg("incoming_request")
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
