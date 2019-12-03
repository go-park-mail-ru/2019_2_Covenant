package artist

import "2019_2_Covenant/internal/models"

type Usecase interface {
	FindLike(name string, count uint64) ([]*models.Artist, error)
	Store(artist *models.Artist) error
	DeleteByID(id uint64) error
	UpdateByID(id uint64, name string) error
	Fetch(count uint64, offset uint64) ([]*models.Artist, uint64, error)
}
