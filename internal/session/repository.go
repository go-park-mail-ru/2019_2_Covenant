package session

import (
	"2019_2_Covenant/internal/models"
)

/*
 *	Repository interface represents the session's repository contract
 */

type Repository interface {
	Get(value string) (*models.Session, error)
	Store(newSession *models.Session) error
	DeleteByID(id uint64) error
}
