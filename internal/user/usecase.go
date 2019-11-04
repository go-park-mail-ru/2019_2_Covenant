package user

import "2019_2_Covenant/internal/models"

type Usecase interface {
	Fetch(count uint64) ([]*models.User, error)
	GetByID(id uint64) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Store(user *models.User) (*models.User, error)
	UpdateAvatar(id uint64, avatarPath string) (*models.User, error)
	UpdateNickname(id uint64, nickname string) (*models.User, error)
}
