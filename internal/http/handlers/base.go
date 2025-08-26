package handlers

import (
	"net/http"
	"strings"
	"todo/internal/http/handlers/auth"
	"todo/internal/http/handlers/todo"
)

type BaseHandler struct {
	AuthHandler  *auth.AuthHandler
	TasksHandler *todo.TasksHandler
}

func (h *BaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/health" && r.Method == "GET":
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
		}(w, r)
	case r.URL.Path == "/tasks" && r.Method == "GET":
		h.TasksHandler.GetTasks(w, r)

	case r.URL.Path == "/tasks" && r.Method == "POST":
		h.TasksHandler.PostTask(w, r)

	case strings.HasPrefix(r.URL.Path, "/tasks/") && r.Method == "GET":
		h.TasksHandler.GetTask(w, r)

	case strings.HasPrefix(r.URL.Path, "/tasks/") && r.Method == "PATCH":
		h.TasksHandler.PatchTask(w, r)

	case strings.HasPrefix(r.URL.Path, "/tasks/") && r.Method == "DELETE":
		h.TasksHandler.DeleteTask(w, r)

	case r.URL.Path == "/register" && r.Method == "POST":
		h.AuthHandler.Register(w, r)

	case r.URL.Path == "/login" && r.Method == "POST":
		h.AuthHandler.Login(w, r)

	case r.URL.Path == "/logout" && r.Method == "POST":
		h.AuthHandler.Logout(w, r)

	default:
		http.NotFound(w, r)
	}
}
