package task_utils

import (
	"time"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"errors"
	"todo/internal/models"
)

func GetTaskUUID(path string) string {
	extracted, _ := strings.CutPrefix(path, "/tasks/")

	return extracted
}

func ValidString(s string) bool {
	for _, c := range s {
		if !(('!' <= c && c <= '~') || (c == ' ')) {
			return false
		}

	}

	return true
}

var allowedOrderBy = map[string]string{
	"title": "title",
	"priority": "priority",
	"completed": "completed",
	"due_date": "due_date",
	"category": "category",
	"created_at": "created_at",
	"updated_at": "updated_at",
}

func GetDynamicQuery(user_id int, r *http.Request) (string, []interface{}, error) {
	condition_query := " WHERE user_id = $1 AND"
	operation_query := ""
	query := ""
	args := []interface{}{}
	arg_ind := 2

	args = append(args, user_id)

	query_params := GetQueryParams(r)

	if query_params["completed"] != "" {
		completed_str := fmt.Sprintf(" completed = $%d", arg_ind)
		condition_query += completed_str + " AND"

		completed_bool, err := strconv.ParseBool(query_params["completed"])

		if err != nil {
			return "", args, err
		}
		args = append(args, completed_bool)
		arg_ind++
	}

	if query_params["category"] != "" {
		category_str := fmt.Sprintf(" category = $%d", arg_ind)
		condition_query += category_str + " AND"
		args = append(args, query_params["category"])
		arg_ind++
	}

	if query_params["due"] != "" {
		due_str := fmt.Sprintf(" due_date <= $%d", arg_ind)
		condition_query += due_str + " AND"
		args = append(args, query_params["due"])
		arg_ind++
	}

	if query_params["search"] != "" {
		search_str := fmt.Sprintf(" title = $%d", arg_ind)
		condition_query += search_str + " AND"
		args = append(args, query_params["search"])
		arg_ind++
	}

	if query_params["priority"] != "" {
		priority_str := fmt.Sprintf(" priority = $%d", arg_ind)
		condition_query += priority_str + " AND"
		args = append(args, query_params["priority"])
		arg_ind++
	}

	condition_query, _ = strings.CutSuffix(condition_query, " AND")

	if query_params["sort"] != "" {
		sort_param := query_params["sort"]
		
		if _, ok := allowedOrderBy[sort_param]; !ok {
			return "", args, errors.New("invalid sort query param")
		}

		sort_str := fmt.Sprintf(" ORDER BY %s", query_params["sort"])
		operation_query += sort_str
	}

	if query_params["limit"] != "" {
		limit, err := strconv.Atoi(query_params["limit"])

		if err != nil {
			return "", args, err
		}

		limit_str := fmt.Sprintf(" LIMIT $%d", arg_ind)
		operation_query += limit_str
		args = append(args, limit)
		arg_ind++
	}

	if condition_query == " WHERE" {
		condition_query = ""
	}

	query += condition_query + operation_query

	return query, args, nil
}

func GetQueryParams(r *http.Request) map[string]string {
	return map[string]string{
		"completed": r.URL.Query().Get("completed"),
		"category":  r.URL.Query().Get("category"),
		"due":       r.URL.Query().Get("due"),
		"search":    r.URL.Query().Get("search"),
		"sort":      r.URL.Query().Get("sort"),
		"limit":     r.URL.Query().Get("limit"),
		"priority":  r.URL.Query().Get("priority"),
	}
}

func TrimSpace(task *models.NewTask) {
	task.Title = strings.TrimSpace(task.Title)
	task.Due_date = strings.TrimSpace(task.Due_date)
	task.Priority = strings.TrimSpace(task.Priority)
	task.Category = strings.TrimSpace(task.Category)
}

func GetUpdateQuery(user_id int, task_uuid string, update_task models.UpdateTask) (string, []any) {
	update_query := "UPDATE tasks SET "
	args  := []any{}
	arg_ind := 1

	if update_task.Title != nil {
		update_query += fmt.Sprintf("title = $%d, ", arg_ind) 

		args = append(args, *update_task.Title)
		arg_ind++
	}

	if update_task.Due_date != nil {
		update_query += fmt.Sprintf("due_date = $%d, ", arg_ind)
		
		args = append(args, *update_task.Due_date)
		arg_ind++
	}

	if update_task.Priority != nil {	
		update_query += fmt.Sprintf("priority = $%d, ", arg_ind)

		args = append(args, *update_task.Priority)
		arg_ind++
	}

	if update_task.Category != nil {	
		update_query += fmt.Sprintf("category = $%d, ", arg_ind)

		args = append(args, *update_task.Category)
		arg_ind++
	}

	if update_task.Completed != nil {
		update_query += fmt.Sprintf("completed = $%d, ", arg_ind)

		args = append(args, *update_task.Completed)
		arg_ind++
	}

	if update_query == "UPDATE tasks SET " {
		update_query = ""
	} else {
		update_query += fmt.Sprintf("updated_at = $%d WHERE user_id = $%d AND id = $%d", arg_ind, arg_ind + 1, arg_ind + 2)
		args = append(args, time.Now(), user_id, task_uuid)
		arg_ind += 2
	}

	return update_query, args
}
