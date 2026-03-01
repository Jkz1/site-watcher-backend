package repository

import (
	"site-checker-backend/internal/models"

	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	DB *sqlx.DB
}

func (r *UserRepo) Create(u *models.User) error {
	_, err := r.DB.NamedExec(`INSERT INTO users (username, password) VALUES (:username, :password)`, u)
	return err
}

func (r *UserRepo) GetByID(id int) (*models.User, error) {
	var u models.User
	err := r.DB.Get(&u, "SELECT * FROM users WHERE id=$1", id)
	return &u, err
}

func (r *UserRepo) GetByUsername(username string) (*models.User, error) {
	var u models.User
	err := r.DB.Get(&u, "SELECT * FROM users WHERE username=$1", username)
	return &u, err
}

func (r *UserRepo) UpdatePassword(id int, newHash string) error {
	_, err := r.DB.Exec("UPDATE users SET password=$1 WHERE id=$2", newHash, id)
	return err
}
