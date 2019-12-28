package repository

import (
	. "2019_2_Covenant/tools/vars"
	"context"
	"database/sql"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"time"
)

var emp = &empty.Empty{}

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

func (sr *SessionRepository) Get(ctx context.Context, request *GetRequest) (*Session, error) {
	item := &Session{}
	expires := time.Time{}

	if err := sr.db.QueryRow("SELECT id, user_id, expires, data FROM sessions WHERE data = $1",
		request.Value,
	).Scan(
		&item.Id,
		&item.UserId,
		&expires,
		&item.Data,
	); err != nil {
		return nil, err
	}

	if expires.Sub(time.Now()) <= 0 {
		_, err := sr.DeleteByID(ctx, &DeleteByIDRequest{Id: item.Id})

		if err != nil {
			return nil, err
		}

		return nil, ErrExpired
	}

	item.Expires, _ = ptypes.TimestampProto(expires)
	return item, nil
}

func (sr *SessionRepository) Store(_ context.Context, newSession *Session) (*Session, error) {
	expires, _ := ptypes.Timestamp(newSession.Expires)
	if err := sr.db.QueryRow("INSERT INTO sessions (user_id, expires, data) VALUES ($1, $2, $3) RETURNING id",
		newSession.UserId,
		expires,
		newSession.Data,
	).Scan(
		&newSession.Id,
	); err != nil {
		return nil, ErrInternalServerError
	}

	return newSession, nil
}

func (sr *SessionRepository) DeleteByID(_ context.Context, request *DeleteByIDRequest) (*empty.Empty, error) {
	var item uint64

	if err := sr.db.QueryRow("DELETE from sessions WHERE id = $1 RETURNING id",
		request.Id,
	).Scan(
		&item,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, ErrInternalServerError
	}

	return emp, nil
}
