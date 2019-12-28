package playlist

import "2019_2_Covenant/pkg/models"

type Repository interface {
	Store(playlist *models.Playlist) error
	Fetch(userID uint64, count uint64, offset uint64) ([]*models.Playlist, uint64, error)
	DeleteByID(playlistID uint64) error
	AddToPlaylist(playlistID uint64, trackID uint64) error
	RemoveFromPlaylist(playlistID uint64, trackID uint64) error
	GetSinglePlaylist(playlistID uint64) (*models.Playlist, uint64, error)
	GetTracksFrom(playlistID uint64, authID uint64) ([]*models.Track, error)
}
