package handlers

import (
	"net/http"
	"todo/internal/http/handlers/auth"
	"todo/internal/http/handlers/todo"
)

type BaseHandler struct {
	AuthHandler  *auth.AuthHandler
	TasksHandler *todo.TasksHandler
	Mux          *http.ServeMux
}

func (h *BaseHandler) HandleRoutes() {
	h.Mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	h.Mux.HandleFunc("GET /tasks", h.TasksHandler.GetTasks)
	h.Mux.HandleFunc("POST /tasks", h.TasksHandler.PostTask)
	h.Mux.HandleFunc("GET /tasks/{id}", h.TasksHandler.GetTask)
	h.Mux.HandleFunc("PATCH /tasks/{id}", h.TasksHandler.PatchTask)
	h.Mux.HandleFunc("DELETE /tasks/{id}", h.TasksHandler.DeleteTask)
	h.Mux.HandleFunc("POST /register", h.AuthHandler.Register)
	h.Mux.HandleFunc("POST /login", h.AuthHandler.Login)
	h.Mux.HandleFunc("POST /logout", h.AuthHandler.Logout)
}
