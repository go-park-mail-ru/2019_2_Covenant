package user

import (
	"2019_2_Covenant/pkg/models"
)

/*
 *	Repository interface represents the user's repository contract
 */

type Repository interface {
	FetchAll() ([]*models.User, error)
	GetById(id int64) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Store (user *models.User) error
}
