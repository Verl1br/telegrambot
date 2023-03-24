package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CreateTable(cfg Config) error {

	//Connecting to database
	db, err := sqlx.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return err
	}
	defer db.Close()

	//Creating users Table
	if _, err = db.Exec(`CREATE TABLE users(id SERIAL PRIMARY KEY, user_id BIGINT);`); err != nil {
		return err
	}
	if _, err = db.Exec(`CREATE TABLE requests(id SERIAL PRIMARY KEY, user_id INT REFERENCES users (id) ON DELETE CASCADE, text TEXT, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);`); err != nil {
		return err
	}

	return nil
}
