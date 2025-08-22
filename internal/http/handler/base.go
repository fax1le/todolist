package handler

import (
	"todo/internal/http/handler/auth"
	"todo/internal/http/handler/todo"
	"net/http"
	"strings"
)

type BaseHandler struct{}

func (h BaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/health" && r.Method == "GET":
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
		}(w, r)
	case r.URL.Path == "/tasks" && r.Method == "GET":
			tasks.GetTasks(w, r)
	
	case r.URL.Path == "/tasks" && r.Method == "POST":
			tasks.PostTask(w, r)
		
	case strings.HasPrefix(r.URL.Path, "/tasks/") && r.Method == "GET":
			tasks.GetTask(w, r)
		
	case strings.HasPrefix(r.URL.Path, "/tasks/") && r.Method == "PATCH":
			tasks.PatchTask(w, r)
			
	case strings.HasPrefix(r.URL.Path, "/tasks/") && r.Method == "DELETE":
			tasks.DeleteTask(w, r)
		
	case r.URL.Path == "/register" && r.Method == "POST":
			auth.Register(w, r)
	
	case r.URL.Path == "/login" && r.Method == "POST":
			auth.Login(w, r)

	case r.URL.Path == "/logout" && r.Method == "POST":
			auth.Logout(w, r)

	default:
		http.NotFound(w, r)
	}
}
