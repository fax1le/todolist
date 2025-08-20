package session

import (
	"net"
	"net/http"

	"github.com/google/uuid"
)

func MustGenerateUUID() string { return uuid.NewString() }

func SetSessionCookie(w http.ResponseWriter, session_uuid string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session_uuid,
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func GetIP(r *http.Request) string {

	hostPort := r.RemoteAddr
	host,_, err := net.SplitHostPort(hostPort)

	if err != nil {
		return hostPort
	}

	return host
}

func Truncate(s string, max int) string {
	if len(s) <= max { return s }
	return s[:max]
}

