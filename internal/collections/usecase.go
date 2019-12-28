package collections

import "2019_2_Covenant/internal/models"

type Usecase interface {
	Store(collection *models.Collection) error
	DeleteByID(id uint64) error
	UpdateByID(collectionID uint64, name string, description string) error
	Fetch(count uint64, offset uint64) ([]*models.Collection, uint64, error)
	GetByID(id uint64) (*models.Collection, uint64, error)
	AddTrack(collectionID uint64, trackID uint64) error
	GetTracks(collectionID uint64, authID uint64) ([]*models.Track, error)
	UpdatePhoto(collectionID uint64, path string) error
}
