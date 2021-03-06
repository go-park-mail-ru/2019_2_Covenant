package session

import "2019_2_Covenant/internal/models"

type Usecase interface {
	Get(value string) (*models.Session, error)
	Store(newSession *models.Session) error
	DeleteByID(id uint64) error
}
