package repository

import (
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/user"
	. "2019_2_Covenant/tools/vars"
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

func (ur *UserRepository) Store(newUser *models.User) error {
	return ur.db.QueryRow("INSERT INTO users (nickname, email, password) VALUES ($1, $2, $3) RETURNING id, avatar",
		newUser.Nickname,
		newUser.Email,
		newUser.Password,
	).Scan(&newUser.ID, &newUser.Avatar)
}

func (ur *UserRepository) GetByEmail(email string) (*models.User, error) {
	u := &models.User{}

	if err := ur.db.QueryRow("SELECT id, nickname, email, avatar, password, role, access FROM users WHERE email = $1",
		email,
	).Scan(&u.ID, &u.Nickname, &u.Email, &u.Avatar, &u.Password, &u.Role, &u.Access); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return u, nil
}

func (ur *UserRepository) GetByID(usrID uint64) (*models.User, error) {
	u := &models.User{}

	if err := ur.db.QueryRow("SELECT id, nickname, email, avatar, password, role, access FROM users WHERE id = $1",
		usrID,
	).Scan(&u.ID, &u.Nickname, &u.Email, &u.Avatar, &u.Password, &u.Role, &u.Access); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return u, nil
}

func (ur *UserRepository) GetByNickname(nickname string) (*models.User, error) {
	u := &models.User{}

	if err := ur.db.QueryRow("SELECT id, nickname, email, avatar, password, role, access FROM users WHERE nickname = $1",
		nickname,
	).Scan(&u.ID, &u.Nickname, &u.Email, &u.Avatar, &u.Password, &u.Role, &u.Access); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return u, nil
}

func (ur *UserRepository) Fetch(count uint64) ([]*models.User, error) {
	var users []*models.User

	rows, err := ur.db.Query("SELECT id, nickname, email, avatar, password, role, access FROM users LIMIT $1", count)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		u := &models.User{}

		if err := rows.Scan(&u.ID, &u.Nickname, &u.Email, &u.Avatar, &u.Password, &u.Role, &u.Access); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (ur *UserRepository) nicknameExists(nickname string) (bool, error) {
	usr, err := ur.GetByNickname(nickname)

	if err != nil {
		return false, err
	}

	if usr != nil {
		return true, nil
	}

	return false, nil
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

func (ur *UserRepository) emailExists(email string) (bool, error) {
	usr, err := ur.GetByEmail(email)

	if err != nil {
		return false, err
	}

	if usr != nil {
		return true, nil
	}

	return false, nil
}

func (ur *UserRepository) UpdatePassword(id uint64, password string) error {
	if _, err := ur.db.Exec("UPDATE users SET password = $1 WHERE id = $2",
		password,
		id,
	); err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) Update(id uint64, nickname string, email string) (*models.User, error) {
	u := &models.User{}

	if err := ur.db.QueryRow("UPDATE users SET nickname = $1, email = $2 WHERE id = $3 RETURNING nickname, email, avatar",
		nickname,
		email,
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
