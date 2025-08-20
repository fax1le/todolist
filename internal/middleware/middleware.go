package middleware

import (
	"todo/internal/log"
	"todo/internal/storage/redis"
	"todo/internal/http/context"
	"context"
	"net/http"
	"time"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

var PublicRoutes = map[string]string{
	"/health":   "GET",
	"/register": "POST",
	"/login":    "POST",
}


func LoggingMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req_start := time.Now()
		log.Logger.Info("request", "addr", r.RemoteAddr, "path", r.URL, "method", r.Method)

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		log.Logger.Info("request", "response", wrapped.statusCode, "execution time", time.Since(req_start))
	})
}

func AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if val, ok := PublicRoutes[r.URL.Path]; ok && val == r.Method {
			next.ServeHTTP(w, r)
			return
		}

		session_cookie, err := r.Cookie("session_id")

		if err != nil || session_cookie.Value == "" {
			log.Logger.Error("session_id cookie missing", "err", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !redis.SessionExists(session_cookie.Value) {
			log.Logger.Warn("session_id doesn't exist")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user_id, err := redis.GetUID(session_cookie.Value)

		if err != nil {
			log.Logger.Error("Redis error", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), ctx.UserIDKey, user_id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
