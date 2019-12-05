package album

import "2019_2_Covenant/internal/models"

type Usecase interface {
	FindLike(name string, count uint64) ([]*models.Album, error)
	DeleteByID(id uint64) error
	UpdateByID(albumID uint64, artistID uint64, name string, year string) error
	Fetch(count uint64, offset uint64) ([]*models.Album, uint64, error)
	GetByID(id uint64) (*models.Album, uint64, error)
	AddTrack(albumID uint64, track *models.Track) error
}
