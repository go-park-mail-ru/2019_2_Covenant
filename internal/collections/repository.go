package collections

import "2019_2_Covenant/internal/models"

type Repository interface {
	Insert(collection *models.Collection) error
	DeleteByID(id uint64) error
	UpdateByID(collectionID uint64, name string, description string) error
	Select(count uint64, offset uint64) ([]*models.Collection, uint64, error)
	SelectByID(id uint64) (*models.Collection, uint64, error)
	InsertTrack(collectionID uint64, trackID uint64) error
	SelectTracks(collectionID uint64, authID uint64) ([]*models.Track, error)
	UpdatePhoto(collectionID uint64, path string) error
}
