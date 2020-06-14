package storage

type User struct {
	Name      string   `json:"name"`
	ChatID    int64    `json:"chat_id"`
	Interval  int64    `json:"interval"`
	LastReply int64    `json:"last_reply"`
	Symbols   []string `json:"symbols"`
	Suspended bool     `json:"suspended"`
}

type UserStorage interface {
	GetUserByName(name string) (User, error)
	GetNotSuspendedUsers() ([]User, error)
	GetUsersForNotification(now int64) ([]User, error)
	CreateUser(name string, chatID int64) (User, error)
	UpdateUserByName(name string, user User) error
}
