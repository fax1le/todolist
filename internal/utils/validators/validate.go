package validators

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
	"todo/internal/models"
	"todo/internal/storage/postgres"
	"todo/internal/utils/task"
)

const layout = "2006-01-02 15:04:05"

func ValidateTask(user_id int, task models.NewTask) error {
	if task.Title == "" || task.Category == "" {
		return errors.New("insertion requirements not met, can't be empty")
	}

	if !task_utils.ValidString(task.Title) || !task_utils.ValidString(task.Category) {
		return errors.New("insertion requirements not met, not valid string")
	}

	if db.TaskExists(user_id, task.Title) {
		return errors.New("unique task violation: task already exists")
	}

	date, err := time.Parse(layout, task.Due_date)

	if err != nil {
		return errors.New("due time requirements not met, should be YYYY-MM-DD  HH:MM:SS")
	}

	if time.Since(date) >= 0 {
		return errors.New("due time requirements not met, should be > Current time")
	}

	return nil
}

var allowedEmailChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789.-"

func ValidateEmail(email string) error {
	
	if len(email) > 254 {
		return errors.New("email size too large")
	}

	if len(email) < 3 {
		return errors.New("email size too small")
	}

	if strings.Count(email, "@") != 1 {
		return errors.New("@ is missing or more than 1")
	}

	local, domain, _ := strings.Cut(email, "@")

	// local and domain part length validation

	if !(1 <= len(local) && len(local) <= 64) {
		return errors.New("incorrect size of local part")
	}

	if !(1 <= len(domain) && len(domain) <= 253) {
		return errors.New("incorrect size of domain part")
	}

	if strings.Count(domain, ".") == 0 {
		return errors.New("lacking atleast one dot in domain part")
	}

	// Forbidden pattern at the start/end

	start := local[0]
	end := local[len(local)-1]

	if start == '.' || start == '-' || end == '.' || end == '-' {
		return errors.New("cannot start/end with a forbidden char in local part")
	}

	// Splitting domain into labels (abc.def.hij) to identify the start/end of a label

	labels := strings.Split(domain, ".")

	for _, label := range labels {
		if len(label) == 0 {
			return errors.New("cannot start/end with a forbidden char in domain label part")
		}

		start = label[0]
		end = label[len(label)-1]

		if start == '-' || end == '-' {
			return errors.New("cannot start/end with a forbidden char in domain label part")
		}

		// Checking for allowed chars in labels

		for i, char := range label {
			if !strings.ContainsRune(allowedEmailChars, char) {
				return errors.New("forbidden char for domain part")
			}

			if i < len(domain)-1 && (char == '.' || char == '-') && rune(domain[i+1]) == char {
				return errors.New("forbidden consecutive pattern for . or -")
			}
		}
	}

	for i, char := range local {
		if !strings.ContainsRune(allowedEmailChars, char) {
			return errors.New("forbidden char for local part")
		}

		if i < len(local)-1 && (char == '.' || char == '-') && rune(local[i+1]) == char {
			return errors.New("forbidden consecutive pattern for . or - ")
		}

	}

	return nil
}

var allowedPassChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789.-!@#$%^&*()_+={}[]|:;\"'<>,.?/~` "

func ValidatePassword(password string) error {
	if len(password) < 6 {
		return errors.New("password is too short")
	}

	if len(password) > 64 {
		return errors.New("password is too long")
	}

	for _, char := range password {
		if !strings.ContainsRune(allowedPassChars, char) {
			return errors.New("forbidden char for password")
		}

	}

	return nil
}

func GetValidateUpdateParams(user_id int, r *http.Request) (models.UpdateTask, error) {
	// Title

	var update_task models.UpdateTask

	err := json.NewDecoder(r.Body).Decode(&update_task)

	if err != nil {
		return update_task, err
	}

	if update_task.Title != nil {
		if *update_task.Title == "" {
			return update_task, errors.New("update requirements not met, can't be empty")
		}

		if !task_utils.ValidString(*update_task.Title) {
			return update_task, errors.New("update requirements not met, not valid string")
		}

		if db.TaskExists(user_id, *update_task.Title) {
			return update_task, errors.New("unique task violation: task already exists")
		}
	}

	// Due_date

	if update_task.Due_date != nil {
		date, err := time.Parse(layout, *update_task.Due_date)

		if err != nil {
			return update_task, errors.New("due time requirements not met, should be YYYY-MM-DD  HH:MM:SS")
		}

		if time.Since(date) >= 0 {
			return update_task, errors.New("due time requirements not met, should be > Current time")
		}
	}

	// Priority

	if update_task.Priority != nil {
		if *update_task.Priority == "" {
			return update_task, errors.New("update requirements not met, can't be empty")
		}

		if *update_task.Priority != "low" && *update_task.Priority != "medium" && *update_task.Priority != "high" {
			return update_task, errors.New("update requirements not met, priority must be in ('low', 'medium', 'high')")
		}
	}

	// Category

	if update_task.Category != nil {
		if *update_task.Category == "" {
			return update_task, errors.New("update requirements not met, can't be empty")
		}

		if !task_utils.ValidString(*update_task.Category) {
			return update_task, errors.New("update requirements not met, not valid string")
		}
	}

	return update_task, nil
}
