package session

import (
	"github.com/google/uuid"
	"net/http"
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
