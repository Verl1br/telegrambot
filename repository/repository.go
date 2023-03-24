package repository

import (
	"telegram-tz/models"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user int) (int, error)
	GetUser(uesr_id int) (models.User, error)
}

type Request interface {
	CreateRequest(user_id int, text string) (int, error)
	GetRequests(uesr_id int) (int, error)
}

type Repository struct {
	Authorization
	Request
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthRepository(db),
		Request:       NewRequestRepository(db),
	}
}
