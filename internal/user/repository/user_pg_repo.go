package repository

import (
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/user"
	"database/sql"
	"fmt"
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
		fmt.Println(err)
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
		&u.Password,
	); err != nil {
		return nil, err
	}

	return u, nil
}

func (ur *UserRepository) GetByID(usrID uint64) (*models.User, error) {
	u := &models.User{}

	if err := ur.db.QueryRow("SELECT id, nickname, email, password FROM users WHERE id = $1",
		usrID,
	).Scan(
		&u.ID,
		&u.Nickname,
		&u.Email,
		&u.Password,
	); err != nil {
		return nil, err
	}

	return u, nil
}

func (ur *UserRepository) FetchAll(count uint64) ([]*models.User, error) {
	var users []*models.User

	rows, err := ur.db.Query("SELECT id, nickname, email, password FROM users LIMIT $1", count)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			id       uint64
			nickname string
			email    string
			password string
		)

		if err := rows.Scan(&id, &nickname, &email, &password); err != nil {
			return nil, err
		}

		users = append(users, &models.User{
			ID: id,
			Nickname: nickname,
			Email: email,
			Password: password,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (ur *UserRepository) Update(id uint64, name string, surname string) (*models.User, error) {
	u := &models.User{}

	if err := ur.db.QueryRow("UPDATE users SET name = $1, surname = $2 WHERE id = $3 RETURNING nickname, email, name, surname",
		name,
		surname,
		id,
	).Scan(
		&u.Nickname,
		&u.Email,
		&u.Name,
		&u.Surname,
	); err != nil {
		return nil, err
	}

	return u, nil
}
