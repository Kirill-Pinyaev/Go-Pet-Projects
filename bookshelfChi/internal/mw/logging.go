package mw

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func Logging(l *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			reqID := middleware.GetReqID(r.Context())
			l.Printf("[start] id=%s %s %s", reqID, r.Method, r.URL.Path)
			next.ServeHTTP(ww, r)
			l.Printf("[done ] id=%s %s %s status=%d bytes=%d",
				reqID, r.Method, r.URL.Path, ww.Status(), ww.BytesWritten())
		})
	}
}
