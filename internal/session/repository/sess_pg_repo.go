package repository

import (
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/session"
	. "2019_2_Covenant/tools/vars"
	"database/sql"
	"time"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) session.Repository {
	return &SessionRepository{
		db: db,
	}
}

func (sr *SessionRepository) Get(value string) (*models.Session, error) {
	item := &models.Session{}

	if err := sr.db.QueryRow("SELECT id, user_id, expires, data FROM sessions WHERE data = $1",
		value,
	).Scan(
		&item.ID,
		&item.UserID,
		&item.Expires,
		&item.Data,
	); err != nil {
		return nil, err
	}
	timeNow := time.Now()
	diffTime := item.Expires.Sub(timeNow)

	if diffTime <= 0 {
		err := sr.DeleteByID(item.ID)

		if err != nil {
			return nil, err
		}

		return nil, ErrExpired
	}

	return item, nil
}

func (sr *SessionRepository) Store(newSession *models.Session) error {
	if err := sr.db.QueryRow("INSERT INTO sessions (user_id, expires, data) VALUES ($1, $2, $3) RETURNING id",
		newSession.UserID,
		newSession.Expires,
		newSession.Data,
	).Scan(
		&newSession.ID,
	); err != nil {
		return ErrInternalServerError
	}

	return nil
}

func (sr *SessionRepository) DeleteByID(id uint64) error {
	var item uint64

	if err := sr.db.QueryRow("DELETE from sessions WHERE id = $1 RETURNING id",
		id,
	).Scan(
		&item,
	); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}

		return ErrInternalServerError
	}

	return nil
}
