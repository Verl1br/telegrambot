package repository

import (
	"fmt"
	"telegram-tz/models"

	"github.com/jmoiron/sqlx"
)

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(user int) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (user_id) VALUES ($1) RETURNING id", "users")

	row := r.db.QueryRow(query, user)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthRepository) GetUser(user_id int) (models.User, error) {
	var user models.User

	query := fmt.Sprintf("SELECT id, user_id FROM %s WHERE user_id = $1", "users")
	err := r.db.Get(&user, query, user_id)

	return user, err
}
