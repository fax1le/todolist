package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
	"todo/internal/models"
	"todo/internal/storage/postgres"
	redis_ "todo/internal/storage/redis"
	"todo/internal/utils/password"
	"todo/internal/utils/session"
	"todo/internal/utils/validators"

	"github.com/redis/go-redis/v9"
)

type AuthHandler struct {
	DB     *sql.DB
	Cache  *redis.Client
	Logger *slog.Logger
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		h.Logger.Error("Register error", "err", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = validators.ValidateEmail(user.Email)
	if err != nil {
		h.Logger.Error("Email format error", "email", user.Email, "err", err)
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	err = validators.ValidatePassword(user.Password)
	if err != nil {
		h.Logger.Error("Password format error", "err", err)
		http.Error(w, "Invalid password format", http.StatusBadRequest)
		return
	}

	db_ctx, db_cancel := context.WithTimeout(r.Context(), time.Second * 3)
	defer db_cancel()

	if postgres.UserExistsByEmail(h.DB, db_ctx, user.Email) {
		h.Logger.Warn("User with provided email exists", "user", user.Email)
		http.Error(w, "User with provided email exists", http.StatusConflict)
		return
	}

	err = postgres.CreateUser(h.DB, db_ctx, user)
	if err != nil {
		h.Logger.Error("User creation error", "user", user.Email, "err", err)
		http.Error(w, "Failed to register", http.StatusInternalServerError)
		return
	}

	h.Logger.Info("User created: ", "user", user.Email)
	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")

	defer r.Body.Close()

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		h.Logger.Error("Login error", "err", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = validators.ValidateEmail(user.Email)
	if err != nil {
		h.Logger.Error("Email format error", "email", user.Email, "err", err)
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	db_ctx, db_cancel := context.WithTimeout(r.Context(), time.Second * 3)
	defer db_cancel()

	if !postgres.UserExistsByEmail(h.DB, db_ctx, user.Email) {
		h.Logger.Warn("User does not exist", "user", user.Email)
		time.Sleep(time.Second)
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	hashed_password, err := postgres.GetPassword(h.DB, db_ctx, user.Email)
	if err != nil {
		h.Logger.Error("Get password error", "err", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if !password.IsCorrectPassword([]byte(hashed_password), []byte(user.Password)) {
		h.Logger.Warn("Invalid password for user", "user", user.Email)
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	user_id, err := postgres.GetUserID(h.DB, db_ctx, user.Email)
	if err != nil {
		h.Logger.Error("Get user_id error", "err", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	session_uuid := session.MustGenerateUUID()

	ip := session.GetIP(r)

	ua := session.Truncate(r.UserAgent(), 200)

	cache_ctx, cache_cancel := context.WithTimeout(r.Context(), time.Second)
	defer cache_cancel()

	err = redis_.StoreSession(h.Cache, cache_ctx, session_uuid, user_id, ip, ua)
	if err != nil {
		h.Logger.Error("Failed to save refresh token", "user", user.Email, "err", err)
		http.Error(w, "Login failed", http.StatusInternalServerError)
		return
	}

	session.SetSessionCookie(w, session_uuid)

	h.Logger.Info("Logged in", "user", user.Email)
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Vary", "Cookie")

	session_cookie, err := r.Cookie("session_id")
	if err != nil || session_cookie.Value == "" {
		h.Logger.Error("session_id cookie error", "err", err)
		session.ClearSessionCookie(w)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	cache_ctx, cache_cancel := context.WithTimeout(r.Context(), time.Second)
	defer cache_cancel()

	user_session_data, err := redis_.GetDeleteSession(h.Cache, cache_ctx, session_cookie.Value)
	if err != nil {
		h.Logger.Info("session_id not found", "err", err)
	} else {
		h.Logger.Info("User successfully logged out", "user", user_session_data)
	}

	session.ClearSessionCookie(w)

	h.Logger.Info("User logged out")
	w.WriteHeader(http.StatusNoContent)
}
