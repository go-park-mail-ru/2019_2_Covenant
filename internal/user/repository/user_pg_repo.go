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
	if err := ur.db.QueryRow("INSERT INTO users (nickname, email, password) VALUES ($1, $2, $3) RETURNING id, avatar",
		newUser.Nickname,
		newUser.Email,
		newUser.Password,
	).Scan(&newUser.ID, &newUser.Avatar); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (ur *UserRepository) GetByEmail(email string) (*models.User, error) {
	u := &models.User{}

	if err := ur.db.QueryRow("SELECT id, nickname, email, avatar, password FROM users WHERE email = $1",
		email,
	).Scan(&u.ID, &u.Nickname, &u.Email, &u.Avatar, &u.Password); err != nil {
		return nil, err
	}

	return u, nil
}

func (ur *UserRepository) GetByID(usrID uint64) (*models.User, error) {
	u := &models.User{}

	if err := ur.db.QueryRow("SELECT id, nickname, email, avatar, password FROM users WHERE id = $1",
		usrID,
	).Scan(&u.ID, &u.Nickname, &u.Email, &u.Avatar, &u.Password); err != nil {
		return nil, err
	}

	return u, nil
}

func (ur *UserRepository) FetchAll(count uint64) ([]*models.User, error) {
	var users []*models.User

	rows, err := ur.db.Query("SELECT id, nickname, email, avatar, password FROM users LIMIT $1", count)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			id       uint64
			nickname string
			email   string
			avatar    string
			password string
		)

		if err := rows.Scan(&id, &nickname, &email, &avatar, &password); err != nil {
			return nil, err
		}

		users = append(users, &models.User{
			ID: id,
			Nickname: nickname,
			Avatar: avatar,
			Email: email,
			Password: password,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (ur *UserRepository) UpdateAvatar(id uint64, avatarPath string) (*models.User, error) {
	u := &models.User{}

	if err := ur.db.QueryRow("UPDATE users SET avatar = $1 WHERE id = $2 RETURNING nickname, email, avatar",
		avatarPath,
		id,
	).Scan(
		&u.Nickname,
		&u.Email,
		&u.Avatar,
	); err != nil {
		return nil, err
	}

	return u, nil
}
