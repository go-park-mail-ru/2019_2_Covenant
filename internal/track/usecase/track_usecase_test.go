package usecase

import (
	. "2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/track"
	mock "2019_2_Covenant/internal/track/mocks"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
)

//go:generate mockgen -source=../repository.go -destination=../mocks/mock_repository.go -package=mock

type Tracks struct {
	Track []*Track
}

var tracks = Tracks{
	Track: []*Track{
		{ID: 1, Name: "We Are the Champions", Artist: "Queen", Album: "News of the World"},
		{ID: 2, Name: "bad guy", Artist: "Billie Eilish", Album: "WHEN WE ALL FALL ASLEEP, WHERE DO WE GO?"},
		{ID: 3, Name: "Still Loving You", Artist: "Scorpions", Album: "Love at First Sting"},
	},
}

func configTrackUsecase(trackRepo *mock.MockRepository) track.Usecase {
	return NewTrackUsecase(trackRepo)
}

func TestTrackUsecase_Fetch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	trackRepo := mock.NewMockRepository(ctrl)

	exe := func(usecase track.Usecase, count uint64) ([]*Track, error) {
		return usecase.Fetch(count)
	}

	t.Run("Test OK", func(t1 *testing.T) {
		count := uint64(3)
		trackRepo.EXPECT().Fetch(count).Return(tracks.Track, nil)
		usecase := configTrackUsecase(trackRepo)

		expTracks, err := exe(usecase, count)

		if gomock.Not(tracks.Track).Matches(expTracks) || err != nil {
			t1.Fail()
		}
	})

	t.Run("Test with error", func(t2 *testing.T) {
		count := uint64(3)
		trackRepo.EXPECT().Fetch(count).Return(nil, fmt.Errorf("some error"))
		usecase := configTrackUsecase(trackRepo)

		expTracks, err := exe(usecase, count)

		if expTracks != nil || err == nil {
			t2.Fail()
		}
	})
}
