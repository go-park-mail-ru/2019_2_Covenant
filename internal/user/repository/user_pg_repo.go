package repository

import (
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/user"
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) user.Repository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) Store(newUser *models.User) (*models.User, error) {
	if err := ur.db.QueryRow("INSERT INTO users (nickname, email, password) VALUES ($1, $2, $3) RETURNING id",
		newUser.Nickname,
		newUser.Email,
		newUser.Password,
	).Scan(&newUser.ID); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (ur *UserRepository) GetByEmail(email string) (*models.User, error) {
	u := &models.User{}

	if err := ur.db.QueryRow("SELECT id, nickname, email, password FROM users WHERE email = $1",
		email,
	).Scan(
		&u.ID,
		&u.Nickname,
		&u.Email,
	); err != nil {
		return nil, err
	}

	return u, nil
}

func (ur *UserRepository) GetByID(usrID uint64) (*models.User, error) {
	u := &models.User{}

	if err := ur.db.QueryRow("SELECT id, nickname, email, password FROM users WHERE email = $1",
		usrID,
	).Scan(
		&u.ID,
		&u.Nickname,
		&u.Email,
	); err != nil {
		return nil, err
	}

	return u, nil
}

func (ur *UserRepository) FetchAll(count uint64) ([]*models.User, error) {
	users := new([]models.User)

	rows, err := ur.db.QueryRow("s")

}
