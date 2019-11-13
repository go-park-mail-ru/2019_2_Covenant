package track

import "2019_2_Covenant/internal/models"

type Repository interface {
	Fetch(count uint64) ([]*models.Track, error)
	StoreFavourite(userID uint64, trackID uint64) error
	RemoveFavourite(userID uint64, trackID uint64) error
	FetchFavourites(userID uint64, count uint64) ([]*models.Track, error)
}
