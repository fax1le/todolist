package models

type Task struct {
	ID         string `json:"id"`
	User_ID    string `json:"user_id"`
	Title      string `json:"title"`
	Completed  bool   `json:"completed"`
	Due_date   string `json:"due"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
	Priority   string `json:"priority"`
	Category   string `json:"category"`
}

type NewTask struct {
	Title    string `json:"title"`
	Due_date string `json:"due"`
	Priority string `json:"priority"`
	Category string `json:"category"`
}

type UpdateTask struct {
	Title     *string `json:"title"`
	Due_date  *string `json:"due"`
	Priority  *string `json:"priority"`
	Category  *string `json:"category"`
	Completed *bool   `json:"completed"`
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type DBuser struct {
	UID        int    `json:"id"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
}

type Session struct {
	UID int    `json:"uid"`
	IAT int64  `json:"iat"`
	EXP int64  `json:"exp"`
	IP  string `json:"ip"`
	UA  string `json:"ua"`
}
