package logging

import (
	"log/slog"
	"net/http"

	m "foodapp/internal/models"
)

type ResponseWriter struct {
	http.ResponseWriter
	Status  int
	UserCtx m.UserContext
}

func (w *ResponseWriter) WriteHeader(status int) {
	w.Status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *ResponseWriter) Write(body []byte) (int, error) {
	if w.Status == 0 {
		w.Status = http.StatusOK
	}

	return w.ResponseWriter.Write(body)
}

func Log(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &ResponseWriter{ResponseWriter: w}

			next.ServeHTTP(rw, r)

			args := []any{
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.Status,
			}

			user := rw.UserCtx
			if user.UserId != 0 {
				args = append(
					args,
					"userId", user.UserId,
					"userRole", user.Role,
				)
			}

			logger.InfoContext(r.Context(), "http request", args...)
		})
	}
}
