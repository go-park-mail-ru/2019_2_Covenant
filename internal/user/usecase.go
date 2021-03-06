package user

import "2019_2_Covenant/internal/models"

type Usecase interface {
	Fetch(count uint64) ([]*models.User, error)
	GetByID(id uint64) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
    GetByNickname(nickname string, authID uint64) (*models.User, error)
	Store(user *models.User) error
	UpdateAvatar(id uint64, avatarPath string) (*models.User, error)
	Update(id uint64, nickname string, email string) (*models.User, error)
	UpdatePassword(id uint64, plainPassword string) error
	GetFollowers(id uint64, count uint64, offset uint64) ([]*models.User, uint64, error)
	GetFollowing(id uint64, count uint64, offset uint64) ([]*models.User, uint64, error)
	FindLike(name string, count uint64) ([]*models.User, error)
}
