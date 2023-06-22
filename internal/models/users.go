package models

type User struct {
	ID       int    `json:"id,omitempty"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
}

type UserAndTask struct {
	User  User   `json:"user,omitempty"`
	Tasks []Task `json:"tasks,omitempty"`
}

type SearchAndStatus struct {
	Search string
	Status string
}

type LoginRequest struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginResponse struct {
	Session string `json:"session,omitempty"`
}
