package track

import "2019_2_Covenant/pkg/models"

type Usecase interface {
	FetchPopular(count uint64, offset uint64, authID uint64) ([]*models.Track, uint64, error)
	FetchFavourites(userID uint64, count uint64, offset uint64) ([]*models.Track, uint64, error)
	StoreFavourite(userID uint64, trackID uint64) error
	RemoveFavourite(userID uint64, trackID uint64) error
	FindLike(name string, count uint64, authID uint64) ([]*models.Track, error)
}
