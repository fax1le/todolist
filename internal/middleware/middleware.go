package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"
	"todo/internal/http/context"
	redis_ "todo/internal/storage/redis"
	"todo/internal/utils/session"

	"github.com/redis/go-redis/v9"
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
	"/health":   "/health",
	"/register": "/register",
	"/login":    "/login",
}

const renewThreshold = 15 * 60

func LoggingMiddleWare(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req_start := time.Now()
		logger.Info("request", "addr", r.RemoteAddr, "path", r.URL, "method", r.Method)

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		logger.Info("request", "response", wrapped.statusCode, "execution time", time.Since(req_start))
	})
}

func AuthMiddleWare(next http.Handler, logger *slog.Logger, cache *redis.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := PublicRoutes[r.URL.Path]; ok {
			next.ServeHTTP(w, r)
			return
		}

		session_cookie, err := r.Cookie("session_id")

		if err != nil || session_cookie.Value == "" {
			logger.Error("Session_id cookie missing", "err", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		redis_ctx, redis_cancel := context.WithTimeout(r.Context(), time.Second)
		defer redis_cancel()

		session_s, err := redis_.GetSession(cache, redis_ctx, session_cookie.Value)

		if err != nil {
			logger.Warn("Session not found", "err", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if session_s.EXP-time.Now().Unix() < renewThreshold {

			err = redis_.RenewSession(cache, redis_ctx, session_cookie.Value)
			if err != nil {
				logger.Error("Redis failed to renew session", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			session.SetSessionCookie(w, session_cookie.Value)
		}

		ctx := context.WithValue(r.Context(), ctx.UserIDKey, session_s.UID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
