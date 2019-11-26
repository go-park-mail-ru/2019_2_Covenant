package playlist

import "2019_2_Covenant/internal/models"

type Repository interface {
	Store(playlist *models.Playlist) error
}
