package repository

import (
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/user"
	. "2019_2_Covenant/tools/vars"
	"database/sql"
	"github.com/sirupsen/logrus"
	"strings"
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
	err := ur.db.QueryRow("INSERT INTO users (nickname, email, password) VALUES ($1, $2, $3) RETURNING id, avatar",
		newUser.Nickname,
		newUser.Email,
		newUser.Password,
	).Scan(&newUser.ID, &newUser.Avatar)

	if err != nil {
		logrus.Info("DB (store user):", err)
	}

	return err
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

func (ur *UserRepository) GetByNickname(nickname string, authID uint64) (*models.User, error) {
	u := &models.User{}
	s := new(bool)

	if err := ur.db.QueryRow("SELECT id, nickname, email, avatar, role, access, " +
		"id in (select subscribed_to from subscriptions where user_id=$1)" +
		"FROM users WHERE nickname = $2",
		authID,
		nickname,
	).Scan(&u.ID, &u.Nickname, &u.Email, &u.Avatar, &u.Role, &u.Access, s); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, err
	}

	if authID != 0 {
		u.Subscription = s
	}

	return u, nil
}

func (ur *UserRepository) FindLike(name string, count uint64) ([]*models.User, error) {
	var users []*models.User

	rows, err := ur.db.Query("SELECT id, nickname, email, avatar, role, access " +
		"FROM users WHERE lower(nickname) like '%' || $1 || '%' LIMIT $2" ,
		strings.ToLower(name), count,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		u := &models.User{}

		if err := rows.Scan(&u.ID, &u.Nickname, &u.Email, &u.Avatar, &u.Role, &u.Access); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
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

func (ur *UserRepository) GetFollowers(id uint64, count uint64, offset uint64) ([]*models.User, uint64, error) {
	var users []*models.User
	var total uint64

	if err := ur.db.QueryRow("select count(*) from subscriptions where subscribed_to=$1",
		id,
	).Scan(&total); err != nil {
		return nil, total, err
	}

	rows, err := ur.db.Query("SELECT U.id, U.nickname, U.email, U.avatar, U.role, U.access FROM users U " +
		"JOIN subscriptions S ON U.id=S.user_id WHERE S.subscribed_to=$1 " +
		"ORDER BY U.nickname LIMIT $2 OFFSET $3", id, count, offset)

	if err != nil {
		return nil, total, err
	}

	defer rows.Close()

	for rows.Next() {
		u := &models.User{}

		if err := rows.Scan(&u.ID, &u.Nickname, &u.Email, &u.Avatar, &u.Role, &u.Access); err != nil {
			return nil, total, err
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, total, err
	}

	return users, total, nil
}

func (ur *UserRepository) GetFollowing(id uint64, count uint64, offset uint64) ([]*models.User, uint64, error) {
	var users []*models.User
	var total uint64

	if err := ur.db.QueryRow("select count(*) from subscriptions where user_id=$1",
		id,
	).Scan(&total); err != nil {
		return nil, total, err
	}

	rows, err := ur.db.Query("SELECT U.id, U.nickname, U.email, U.avatar, U.role, U.access FROM users U " +
		"JOIN subscriptions S ON U.id=S.subscribed_to WHERE S.user_id=$1 " +
		"ORDER BY U.nickname LIMIT $2 OFFSET $3", id, count, offset)

	if err != nil {
		return nil, total, err
	}

	defer rows.Close()

	for rows.Next() {
		u := &models.User{}

		if err := rows.Scan(&u.ID, &u.Nickname, &u.Email, &u.Avatar, &u.Role, &u.Access); err != nil {
			return nil, total, err
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, total, err
	}

	return users, total, nil
}
