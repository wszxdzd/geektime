package domain

import "time"

type User struct {
	Id       int64
	Email    string
	Phone    string
	Password string
	Nickname string
	AboutMe  string
	Birthday time.Time
	Ctime    time.Time
}
