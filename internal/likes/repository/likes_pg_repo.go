package repository

import (
	"2019_2_Covenant/internal/likes"
	. "2019_2_Covenant/tools/vars"
	"database/sql"
	"fmt"
)

type LikesRepository struct {
	db *sql.DB
}

func NewLikesRepository(db *sql.DB) likes.Repository {
	return &LikesRepository{
		db: db,
	}
}

func (ssR *LikesRepository) Like(userID uint64, trackID uint64) error {
	if err := ssR.db.QueryRow("SELECT id FROM likes WHERE user_id = $1 AND track_id = $2",
		userID,
		trackID,
	).Scan(); err == nil {
		return ErrAlreadyExist
	}

	if _, err := ssR.db.Exec("INSERT INTO likes (user_id, track_id) VALUES ($1, $2)",
		userID,
		trackID,
	); err != nil {
		return err
	}

	return nil
}

func (ssR *LikesRepository) Unlike(userID uint64, trackID uint64) error {
	fmt.Println(userID, trackID)
	res, err := ssR.db.Exec("DELETE FROM likes WHERE user_id = $1 AND track_id = $2",
		userID,
		trackID,
	)

	if err != nil {
		return err
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrNotFound
	}

	return nil
}
