package user

import (
	"2019_2_Covenant/pkg/models"
	"context"
	"io"
)

type Usecase interface {
	Fetch(ctx context.Context, count uint64) ([]*models.User, error)
	GetByID(ctx context.Context, id uint64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
    GetByNickname(ctx context.Context, nickname string, authID uint64) (*models.User, error)
	Store(ctx context.Context, user *models.User) error
	Update(ctx context.Context, id uint64, nickname string, email string) (*models.User, error)
	UpdateAvatar(ctx context.Context, userID uint64, photo io.Reader) error
	UpdatePassword(ctx context.Context, id uint64, plainPassword string) error
	GetFollowers(ctx context.Context, id uint64, count uint64, offset uint64) ([]*models.User, uint64, error)
	GetFollowing(ctx context.Context, id uint64, count uint64, offset uint64) ([]*models.User, uint64, error)
	FindLike(ctx context.Context, name string, count uint64) ([]*models.User, error)
}
