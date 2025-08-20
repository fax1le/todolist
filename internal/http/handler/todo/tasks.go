package tasks

import (
	"todo/internal/log"
	"todo/internal/models"
	"todo/internal/storage/postgres"
	"todo/internal/utils/task"
	"todo/internal/utils/validators"
	"todo/internal/http/context"
	"net/http"
	"database/sql"
	"encoding/json"
)

func GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	val := r.Context().Value(ctx.UserIDKey)

	user_id, ok := val.(int)

	if !ok {
		log.Logger.Error("request: failed to get context key value")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query_params, args, err := task_utils.GetDynamicQuery(user_id, r)

	if err != nil {
		log.Logger.Error("postgres: dynamic query error", "err", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	tasks, err := db.SelectTasks(query_params, args)

	if err != nil {
		log.Logger.Info("postgres: select tasks error", "err", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tasks)
	w.WriteHeader(http.StatusOK)
}

func PostTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	val := r.Context().Value(ctx.UserIDKey)

	user_id, ok := val.(int)

	if !ok {
		log.Logger.Error("request: failed to get context key value")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var new_task models.NewTask

	err := json.NewDecoder(r.Body).Decode(&new_task)

	if err != nil {
		log.Logger.Error("request: parsing error", "err", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = validators.ValidateTask(user_id, new_task)

	if err != nil {
		log.Logger.Error("validate: task validation error", "err", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = db.InsertTask(user_id, new_task)

	if err != nil {
		log.Logger.Error("postgres: insertion error", "err", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	log.Logger.Info("Task was created", "task", new_task.Title)
	w.WriteHeader(http.StatusCreated)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	val := r.Context().Value(ctx.UserIDKey)

	user_id, ok := val.(int)

	if !ok {
		log.Logger.Error("request: failed to get context key value")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	task_uuid := task_utils.GetTaskUUID(r.URL.Path)

	task, err := db.SelectTask(user_id, task_uuid)

	if err != nil {
		log.Logger.Warn("postgres: task was not found", "err", err)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(task)
	w.WriteHeader(http.StatusOK)
}

func PatchTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	val := r.Context().Value(ctx.UserIDKey)

	user_id, ok := val.(int)

	if !ok {
		log.Logger.Error("request: failed to get context key value")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	task_uuid := task_utils.GetTaskUUID(r.URL.Path)

	err := db.UpdateTask(user_id, task_uuid)

	if err == sql.ErrNoRows {
		log.Logger.Warn("postgres: task was not found", "err", err)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if err != nil {
		log.Logger.Error("postgres: update task error", "err", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	log.Logger.Info("Task was updated")
	w.WriteHeader(http.StatusOK)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	val := r.Context().Value(ctx.UserIDKey)

	user_id, ok := val.(int)

	if !ok {
		log.Logger.Error("request: failed to get context key value")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	task_uuid := task_utils.GetTaskUUID(r.URL.Path)

	rows_affected, err := db.RemoveTask(user_id, task_uuid)

	if rows_affected == 0 {
		log.Logger.Warn("postgres: task was not found", "err", "psql: 0 rows affected")
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if err != nil {
		log.Logger.Error("postgres: delete task error", "err", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	log.Logger.Info("Task was deleted")
	w.WriteHeader(http.StatusOK)
}
