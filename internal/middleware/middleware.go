package middleware

import (
	"context"
	"net/http"
	"time"
	"todo/internal/http/context"
	"todo/internal/log"
	"todo/internal/storage/redis"
	"todo/internal/utils/session"
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

const renewThreshold = 15 * 60

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
			log.Logger.Error("Session_id cookie missing", "err", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		session_s, err := redis.GetSession(session_cookie.Value)

		// improve error handling
		if err != nil {
			log.Logger.Warn("Session not found", "err", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if session_s.EXP-time.Now().Unix() < renewThreshold {

			err = redis.RenewSession(session_cookie.Value)
			if err != nil {
				log.Logger.Error("Redis failed to renew session", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			session.SetSessionCookie(w, session_cookie.Value)
		}

		ctx := context.WithValue(r.Context(), ctx.UserIDKey, session_s.UID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
