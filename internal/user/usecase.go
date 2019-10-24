package user

import "2019_2_Covenant/internal/models"

type Usecase interface {
	FetchAll() ([]*models.User, error)
	GetByID(id uint64) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Store(user *models.User) error
}
