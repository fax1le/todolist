package todo

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
	"todo/internal/http/context"
	"todo/internal/models"
	"todo/internal/storage/postgres"
	"todo/internal/utils/task"
	"todo/internal/utils/validators"

	"github.com/redis/go-redis/v9"
)


type TasksHandler struct {
	DB *sql.DB
	Cache *redis.Client
	Logger *slog.Logger
}


func (h *TasksHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	val := r.Context().Value(ctx.UserIDKey)

	user_id, ok := val.(int)
	if !ok {
		h.Logger.Error("request: failed to get context key value")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query_params, args, err := task_utils.GetDynamicQuery(user_id, r)
	if err != nil {
		h.Logger.Error("dynamic query error", "err", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	db_ctx, db_cancel := context.WithTimeout(r.Context(), time.Second * 3)
	defer db_cancel()

	tasks, err := postgres.SelectTasks(h.DB, db_ctx, query_params, args)
	if err != nil {
		h.Logger.Info("postgres: select tasks error", "err", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

func (h *TasksHandler) PostTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	val := r.Context().Value(ctx.UserIDKey)

	user_id, ok := val.(int)
	if !ok {
		h.Logger.Error("request: failed to get context key value")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var new_task models.NewTask

	err := json.NewDecoder(r.Body).Decode(&new_task)
	if err != nil {
		h.Logger.Error("request: parsing error", "err", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	db_ctx, db_cancel := context.WithTimeout(r.Context(), time.Second * 3)
	defer db_cancel()

	err = validators.ValidateTask(h.DB, db_ctx, user_id, new_task)
	if err != nil {
		h.Logger.Error("validate: task validation error", "err", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = postgres.InsertTask(h.DB, db_ctx, user_id, new_task)
	if err != nil {
		h.Logger.Error("postgres: insertion error", "err", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Task was created", "task", new_task.Title)
	w.WriteHeader(http.StatusCreated)
}

func (h *TasksHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	val := r.Context().Value(ctx.UserIDKey)

	user_id, ok := val.(int)
	if !ok {
		h.Logger.Error("request: failed to get context key value")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	task_uuid := task_utils.GetTaskUUID(r.URL.Path)

	db_ctx, db_cancel := context.WithTimeout(r.Context(), time.Second * 3)
	defer db_cancel()

	task, err := postgres.SelectTask(h.DB, db_ctx, user_id, task_uuid)
	if err != nil {
		h.Logger.Warn("postgres: task was not found", "err", err)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (h *TasksHandler) PatchTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	val := r.Context().Value(ctx.UserIDKey)

	user_id, ok := val.(int)
	if !ok {
		h.Logger.Error("request: failed to get context key value")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	task_uuid := task_utils.GetTaskUUID(r.URL.Path)

	db_ctx, db_cancel := context.WithTimeout(r.Context(), time.Second * 3)
	defer db_cancel()

	update_task, err := validators.GetValidateUpdateParams(h.DB, db_ctx, user_id, r)
	if err != nil {
		h.Logger.Error("validate: update params validation failed", "err", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	update_query, args := task_utils.GetUpdateQuery(user_id, task_uuid, update_task)

	err = postgres.UpdateTask(h.DB, db_ctx, update_query, args)
	if err != nil {
		if err == sql.ErrNoRows {
		h.Logger.Warn("postgres: task was not found", "err", err)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

		h.Logger.Error("postgres: update task error", "err", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Task was updated")
	w.WriteHeader(http.StatusOK)
}

func (h *TasksHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	val := r.Context().Value(ctx.UserIDKey)

	user_id, ok := val.(int)
	if !ok {
		h.Logger.Error("request: failed to get context key value")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	task_uuid := task_utils.GetTaskUUID(r.URL.Path)

	db_ctx, db_cancel := context.WithTimeout(r.Context(), time.Second * 3)
	defer db_cancel()

	rows_affected, err := postgres.RemoveTask(h.DB, db_ctx, user_id, task_uuid)

	if rows_affected == 0 {
		h.Logger.Warn("postgres: task was not found", "err", "psql: 0 rows affected")
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if err != nil {
		h.Logger.Error("postgres: delete task error", "err", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Task was deleted")
	w.WriteHeader(http.StatusOK)
}
