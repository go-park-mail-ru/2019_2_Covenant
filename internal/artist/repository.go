package artist

import "2019_2_Covenant/internal/models"

type Repository interface {
	FindLike(name string, count uint64) ([]*models.Artist, error)
	Store(artist *models.Artist) error
	CreateAlbum(album *models.Album) error
	DeleteByID(id uint64) error
	UpdateByID(id uint64, name string) error
	Fetch(count uint64, offset uint64) ([]*models.Artist, uint64, error)
	GetByID(id uint64) (*models.Artist, uint64, error)
	UpdatePhoto(artistID uint64, path string) error
	GetArtistAlbums(artistID uint64, count uint64, offset uint64) ([]*models.Album, uint64, error)
}
