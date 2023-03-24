package models

import "time"

type Request struct {
	Id     int       `db:"id"`
	UserId int       `db:"user_id"`
	Text   string    `db:"text"`
	Date   time.Time `db:"created_at"`
}
