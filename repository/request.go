package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type RequestRepository struct {
	db *sqlx.DB
}

func NewRequestRepository(db *sqlx.DB) *RequestRepository {
	return &RequestRepository{db: db}
}

func (r *RequestRepository) CreateRequest(user_id int, text string) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (user_id, text) VALUES ($1, $2) RETURNING id", "requests")

	row := r.db.QueryRow(query, user_id, text)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *RequestRepository) GetRequests(user_id int) (int, error) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE user_id = $1", "requests")
	err := r.db.Get(&count, query, user_id)
	return count, err
}
