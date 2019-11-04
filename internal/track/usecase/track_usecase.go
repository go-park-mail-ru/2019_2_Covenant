package usecase

import (
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/track"
)

type trackUsecase struct {
	trackRepo track.Repository
}

func NewTrackUsecase(tr track.Repository) track.Usecase {
	return &trackUsecase{
		trackRepo: tr,
	}
}

func (tUC *trackUsecase) Fetch(count uint64) ([]*models.Track, error) {
	users, err := tUC.trackRepo.Fetch(count)

	if err != nil {
		return nil, err
	}

	return users, nil
}
