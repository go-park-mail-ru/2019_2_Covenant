package usecase

import (
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/track"
	"2019_2_Covenant/tools/time_parser"
)

type trackUsecase struct {
	trackRepo track.Repository
}

func NewTrackUsecase(tr track.Repository) track.Usecase {
	return &trackUsecase{
		trackRepo: tr,
	}
}

func (tUC *trackUsecase) FetchPopular(count uint64, offset uint64, authID uint64) ([]*models.Track, uint64, error) {
	tracks, total, err := tUC.trackRepo.Fetch(count, offset, authID)

	if err != nil {
		return nil, total, err
	}

	if tracks == nil {
		tracks = []*models.Track{}
	}

	for _, item := range tracks { item.Duration = time_parser.GetDuration(item.Duration) }

	return tracks, total, nil
}

func (tUC *trackUsecase) StoreFavourite(userID uint64, trackID uint64) error {
	if err := tUC.trackRepo.StoreFavourite(userID, trackID); err != nil {
		return err
	}

	return nil
}

func (tUC *trackUsecase) RemoveFavourite(userID uint64, trackID uint64) error {
	if err := tUC.trackRepo.RemoveFavourite(userID, trackID); err != nil {
		return err
	}

	return nil
}

func (tUC *trackUsecase) FetchFavourites(userID uint64, count uint64, offset uint64) ([]*models.Track, uint64, error) {
	tracks, total, err := tUC.trackRepo.FetchFavourites(userID, count, offset)

	if err != nil {
		return nil, total, err
	}

	if tracks == nil {
		tracks = []*models.Track{}
	}

	for _, item := range tracks { item.Duration = time_parser.GetDuration(item.Duration) }

	return tracks, total, nil
}

func (tUC *trackUsecase) FindLike(name string, count uint64, authID uint64) ([]*models.Track, error) {
	tracks, err := tUC.trackRepo.FindLike(name, count, authID)

	if err != nil {
		return nil, err
	}

	if tracks == nil {
		tracks = []*models.Track{}
	}

	for _, item := range tracks { item.Duration = time_parser.GetDuration(item.Duration) }

	return tracks, nil
}
