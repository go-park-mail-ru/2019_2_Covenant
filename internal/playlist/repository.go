package playlist

import "2019_2_Covenant/internal/models"

type Repository interface {
	Store(playlist *models.Playlist) error
	Fetch(userID uint64, count uint64, offset uint64) ([]*models.Playlist, uint64, error)
}
