package session

import (
	"2019_2_Covenant/pkg/models"
	"context"
)

type Usecase interface {
	Get(ctx context.Context, value string) (*models.Session, error)
	Store(ctx context.Context, newSession *models.Session) error
	DeleteByID(ctx context.Context, id uint64) error
}
