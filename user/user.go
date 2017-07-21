package user

import "time"

type Model struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Password  string
}

func (Model) TableName() string {
	return "user"
}

type Session struct {
	ID        int
	UserID    int
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Session) TableName() string {
	return "session"
}
