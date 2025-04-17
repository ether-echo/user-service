package domain

type User struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
	ChatId    int64  `json:"chat_id,omitempty"`
	Message   string `json:"message,omitempty"`
}
