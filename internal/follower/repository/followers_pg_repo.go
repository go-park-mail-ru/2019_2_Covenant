package repository

import (
	"2019_2_Covenant/internal/follower"
	. "2019_2_Covenant/tools/vars"
	"database/sql"
)

type FollowerRepository struct {
	db *sql.DB
}

func NewFollowerRepository(db *sql.DB) follower.Repository {
	return &FollowerRepository{
		db: db,
	}
}

func (flR *FollowerRepository) Subscribe(userID uint64, followerId uint64) error {
	if err := flR.db.QueryRow("SELECT id FROM follower WHERE user_id = $1 AND follower_id = $2",
		userID,
		followerId,
	).Scan(); err == nil {
		return ErrAlreadyExist
	}

	if _, err := flR.db.Exec("INSERT INTO follower (user_id, follower_id) VALUES ($1, $2)",
		userID,
		followerId,
	); err != nil {
			return err
	}

	return nil
}

func (flR *FollowerRepository) Unsubscribe(userID uint64, followerId uint64) error {
	res, err := flR.db.Exec("DELETE FROM follower WHERE user_id = $1 AND follower_id = $2",
		userID,
		followerId,
	)

	if err != nil {
		return err
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrNotFound
	}

	return nil
}
