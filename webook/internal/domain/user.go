package domain

import "time"

type User struct {
	Id       int64
	Nickname string
	Email    string
	Password string
	Phone    string
	About    string
	Birthday string
	Ctime    time.Time
}
