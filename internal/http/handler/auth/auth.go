package auth

import (
	"todo/internal/log"
	"todo/internal/models"
	"todo/internal/storage/postgres"
	"todo/internal/storage/redis"
	"todo/internal/utils/password"
	"todo/internal/utils/session"
	"todo/internal/utils/validators"
	"net/http"
	"encoding/json"
	"time"
)

func Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Logger.Error("Register error", "err", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = validators.ValidateEmail(user.Email)
	if err != nil {
		log.Logger.Error("Email format error", "email", user.Email, "err", err)
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	err = validators.ValidatePassword(user.Password)
	if err != nil {
		log.Logger.Error("Password format error", "password", user.Password, "err", err)
		http.Error(w, "Invalid password format", http.StatusBadRequest)
		return
	}

	if db.UserExistsByEmail(user.Email) {
		log.Logger.Warn("User with provided email exists", "user", user.Email)
		http.Error(w, "User with provided email exists", http.StatusConflict)
		return
	}

	err = db.CreateUser(user)
	if err != nil {
		log.Logger.Error("User creation error", "user", user.Email, "err", err)
		http.Error(w, "Failed to register", http.StatusInternalServerError)
		return
	}

	log.Logger.Info("User created: ", "user", user.Email)
	w.WriteHeader(http.StatusCreated)
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")

	defer r.Body.Close()

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Logger.Error("Login error", "err", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = validators.ValidateEmail(user.Email)
	if err != nil {
		log.Logger.Error("Email format error", "email", user.Email, "err", err)
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	if !db.UserExistsByEmail(user.Email) {
		log.Logger.Warn("User does not exist", "user", user.Email)

		time.Sleep(time.Second)
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	hashed_password := db.MustGetPassword(user.Email)

	if !password.IsCorrectPassword([]byte(hashed_password), []byte(user.Password)) {
		log.Logger.Warn("Invalid password for user", "user", user.Email)
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	user_id := db.MustGetUserID(user.Email)

	session_uuid := session.MustGenerateUUID()

	err = redis.StoreSession(session_uuid, user_id)
	if err != nil {
		log.Logger.Error("Failed to save refresh token", "user", user.Email, "err", err)
		http.Error(w, "Login failed", http.StatusInternalServerError)
		return
	}

	session.SetSessionCookie(w, session_uuid)

	log.Logger.Info("Logged in", "user", user.Email)
	w.WriteHeader(http.StatusOK)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Vary", "Cookie")

	session_cookie, err := r.Cookie("session_id")
	if err != nil || session_cookie.Value == "" {
		log.Logger.Error("session_id cookie error", "err", err)
		session.ClearSessionCookie(w)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	idStr, err := redis.GetDeleteSession(session_cookie.Value)
	if err != nil {
		log.Logger.Info("session_id not found", "err", err)
	} else {
		log.Logger.Info("User successfully logged out", "user", idStr)
	}

	session.ClearSessionCookie(w)

	log.Logger.Info("User logged out")
	w.WriteHeader(http.StatusNoContent)
}
