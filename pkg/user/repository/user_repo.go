package repository

import (
	. "2019_2_Covenant/tools/vars"
	"context"
	"database/sql"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"strings"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) Fetch(_ context.Context, request *FetchRequest) (*UserArray, error) {
	userArray := &UserArray{}

	rows, err := ur.db.Query("SELECT id, nickname, email, avatar, password, role, access FROM users LIMIT $1", request.Count)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		u := &User{}

		if err := rows.Scan(&u.Id, &u.Nickname, &u.Email, &u.Avatar, &u.Password, &u.Role, &u.Access); err != nil {
			return nil, err
		}

		userArray.Users = append(userArray.Users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userArray, nil
}

func (ur *UserRepository) GetByID(_ context.Context, request *GetByIDRequest) (*User, error) {
	u := &User{}

	if err := ur.db.QueryRow("SELECT id, nickname, email, avatar, password, role, access FROM users WHERE id = $1",
		request.Id,
	).Scan(&u.Id, &u.Nickname, &u.Email, &u.Avatar, &u.Password, &u.Role, &u.Access); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return u, nil
}

func (ur *UserRepository) GetByEmail(_ context.Context, request *GetByEmailRequest) (*User, error) {
	u := &User{}

	if err := ur.db.QueryRow("SELECT id, nickname, email, avatar, password, role, access FROM users WHERE email = $1",
		request.Email,
	).Scan(&u.Id, &u.Nickname, &u.Email, &u.Avatar, &u.Password, &u.Role, &u.Access); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return u, nil
}

func (ur *UserRepository) GetByNickname(_ context.Context, request *GetByNicknameRequest) (*User, error) {
	u := &User{}
	var s bool

	if err := ur.db.QueryRow("SELECT id, nickname, email, avatar, role, access, "+
		"id in (select subscribed_to from subscriptions where user_id=$1)"+
		"FROM users WHERE nickname = $2",
		request.AuthID,
		request.Nickname,
	).Scan(&u.Id, &u.Nickname, &u.Email, &u.Avatar, &u.Role, &u.Access, &s); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, err
	}

	if request.AuthID != 0 {
		u.Subscription = s
	}

	return u, nil
}

func (ur *UserRepository) Store(_ context.Context, newUser *User) (*User, error) {
	err := ur.db.QueryRow("INSERT INTO users (nickname, email, password) VALUES ($1, $2, $3) RETURNING id, avatar",
		newUser.Nickname,
		newUser.Email,
		newUser.Password,
	).Scan(&newUser.Id, &newUser.Avatar)

	if err != nil {
		logrus.Info("DB (store user):", err)
		return nil, err
	}

	return newUser, nil
}

func (ur *UserRepository) UpdatePassword(_ context.Context, request *UpdatePasswordRequest) (*empty.Empty, error) {
	if _, err := ur.db.Exec("UPDATE users SET password = $1 WHERE id = $2",
		request.Password,
		request.Id,
	); err != nil {
		return nil, err
	}

	return new(empty.Empty), nil
}

func (ur *UserRepository) Update(_ context.Context, request *UpdateRequest) (*User, error) {
	u := &User{}

	if err := ur.db.QueryRow("UPDATE users SET nickname = $1, email = $2 WHERE id = $3 RETURNING nickname, email, avatar",
		request.Nickname,
		request.Email,
		request.Id,
	).Scan(
		&u.Nickname,
		&u.Email,
		&u.Avatar,
	); err != nil {
		return nil, err
	}

	return u, nil
}

func (ur *UserRepository) GetFollowers(_ context.Context, request *GetFollowRequest) (*GetFollowResponse, error) {
	var total uint64

	if err := ur.db.QueryRow("select count(*) from subscriptions where subscribed_to=$1",
		request.Id,
	).Scan(&total); err != nil {
		return nil, err
	}

	rows, err := ur.db.Query("SELECT U.id, U.nickname, U.email, U.avatar, U.role, U.access FROM users U "+
		"JOIN subscriptions S ON U.id=S.user_id WHERE S.subscribed_to=$1 "+
		"ORDER BY U.nickname LIMIT $2 OFFSET $3", request.Id, request.Count, request.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var users []*User

	for rows.Next() {
		u := &User{}

		if err := rows.Scan(&u.Id, &u.Nickname, &u.Email, &u.Avatar, &u.Role, &u.Access); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &GetFollowResponse{Users: users, Total: total}, nil
}

func (ur *UserRepository) GetFollowing(_ context.Context, request *GetFollowRequest) (*GetFollowResponse, error) {
	var total uint64

	if err := ur.db.QueryRow("select count(*) from subscriptions where user_id=$1",
		request.Id,
	).Scan(&total); err != nil {
		return nil, err
	}

	rows, err := ur.db.Query("SELECT U.id, U.nickname, U.email, U.avatar, U.role, U.access FROM users U "+
		"JOIN subscriptions S ON U.id=S.subscribed_to WHERE S.user_id=$1 "+
		"ORDER BY U.nickname LIMIT $2 OFFSET $3", request.Id, request.Count, request.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var users []*User

	for rows.Next() {
		u := &User{}

		if err := rows.Scan(&u.Id, &u.Nickname, &u.Email, &u.Avatar, &u.Role, &u.Access); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &GetFollowResponse{Users: users, Total: total}, nil
}

func (ur *UserRepository) FindLike(_ context.Context, request *FindLikeRequest) (*UserArray, error) {
	userArray := &UserArray{}

	rows, err := ur.db.Query("SELECT id, nickname, email, avatar, role, access "+
		"FROM users WHERE lower(nickname) like '%' || $1 || '%' LIMIT $2",
		strings.ToLower(request.Name), request.Count,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		u := &User{}

		if err := rows.Scan(&u.Id, &u.Nickname, &u.Email, &u.Avatar, &u.Role, &u.Access); err != nil {
			return nil, err
		}

		userArray.Users = append(userArray.Users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userArray, nil
}

func (ur *UserRepository) emailExists(ctx context.Context, email string) (bool, error) {
	usr, err := ur.GetByEmail(ctx, &GetByEmailRequest{Email: email})

	if err != nil {
		return false, err
	}

	if usr != nil {
		return true, nil
	}

	return false, nil
}
