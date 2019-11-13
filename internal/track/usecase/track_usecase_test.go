package usecase

import (
	. "2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/track"
	mock "2019_2_Covenant/internal/track/mocks"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
)

// Для тестирования только этого файла:
// go test -v -cover -race ./internal/track/usecase

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

func TestTrackUsecase_FetchPopular(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	trackRepo := mock.NewMockRepository(ctrl)

	exe := func(usecase track.Usecase, count uint64) ([]*Track, error) {
		return usecase.FetchPopular(count)
	}

	t.Run("Test OK", func(t1 *testing.T) {
		count := uint64(3)
		trackRepo.EXPECT().Fetch(count).Return(tracks.Track, nil)
		usecase := configTrackUsecase(trackRepo)

		expTracks, err := exe(usecase, count)

		if gomock.Not(tracks.Track).Matches(expTracks) || err != nil {
			fmt.Println("Error. Expected: nil, Got:", err)
			fmt.Println("expTracks. Expected:", tracks.Track, "Got:", expTracks)
			t1.Fail()
		}
	})

	t.Run("Error tracks not exist", func(t2 *testing.T) {
		count := uint64(3)
		trackRepo.EXPECT().Fetch(count).Return(nil, nil)
		usecase := configTrackUsecase(trackRepo)

		expTracks, err := exe(usecase, count)

		if gomock.Not([]*Track{}).Matches(expTracks) || err != nil {
			fmt.Println("Error. Expected: nil, Got: ", err)
			fmt.Println("expTracks. Expected:", []Tracks{} , "Got:", expTracks)
			t2.Fail()
		}
	})

	t.Run("Error fetching", func(t3 *testing.T) {
		count := uint64(3)
		trackRepo.EXPECT().Fetch(count).Return(nil, fmt.Errorf("some error"))
		usecase := configTrackUsecase(trackRepo)

		expTracks, err := exe(usecase, count)

		if expTracks != nil || err == nil {
			fmt.Println("Error. Expected: not nil, Got: ", err)
			fmt.Println("expTracks. Expected: nil, Got: ", expTracks)
			t3.Fail()
		}
	})
}

func TestTrackUsecase_StoreFavourite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	trackRepo := mock.NewMockRepository(ctrl)

	exe := func(usecase track.Usecase, userID uint64, trackID uint64) error {
		return usecase.StoreFavourite(userID, trackID)
	}

	t.Run("Test OK", func(t1 *testing.T) {
		trackRepo.EXPECT().StoreFavourite(gomock.Any(), gomock.Any()).Return(nil)
		usecase := configTrackUsecase(trackRepo)

		err := exe(usecase, uint64(1), uint64(1))

		if err != nil {
			fmt.Println("error happens: ", err)
			t1.Fail()
		}
	})

	t.Run("Test with error", func(t2 *testing.T) {
		trackRepo.EXPECT().StoreFavourite(gomock.Any(), gomock.Any()).Return(fmt.Errorf("some error"))
		usecase := configTrackUsecase(trackRepo)

		err := exe(usecase, uint64(1), uint64(1))

		if err == nil {
			fmt.Println("error happens: ", err)
			t2.Fail()
		}
	})
}

func TestTrackUsecase_RemoveFavourite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	trackRepo := mock.NewMockRepository(ctrl)

	exe := func(usecase track.Usecase, userID uint64, trackID uint64) error {
		return usecase.RemoveFavourite(userID, trackID)
	}

	t.Run("Test OK", func(t1 *testing.T) {
		trackRepo.EXPECT().RemoveFavourite(gomock.Any(), gomock.Any()).Return(nil)
		usecase := configTrackUsecase(trackRepo)

		err := exe(usecase, uint64(1), uint64(1))

		if err != nil {
			fmt.Println("error happens: ", err)
			t1.Fail()
		}
	})

	t.Run("Test with error", func(t2 *testing.T) {
		trackRepo.EXPECT().RemoveFavourite(gomock.Any(), gomock.Any()).Return(fmt.Errorf("some error"))
		usecase := configTrackUsecase(trackRepo)

		err := exe(usecase, uint64(1), uint64(1))

		if err == nil {
			fmt.Println("error happens: ", err)
			t2.Fail()
		}
	})
}

func TestTrackUsecase_FetchFavourites(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	trackRepo := mock.NewMockRepository(ctrl)

	exe := func(usecase track.Usecase, userID uint64, count uint64) ([]*Track, error) {
		return usecase.FetchFavourites(userID, count)
	}

	t.Run("Test OK", func(t1 *testing.T) {
		count := uint64(3)
		userID := uint64(1)
		trackRepo.EXPECT().FetchFavourites(userID, count).Return(tracks.Track, nil)
		usecase := configTrackUsecase(trackRepo)

		expTracks, err := exe(usecase, userID, count)

		if gomock.Not(tracks.Track).Matches(expTracks) || err != nil {
			fmt.Println("Error. Expected: nil, Got:", err)
			fmt.Println("expTracks. Expected:", tracks.Track, "Got:", expTracks)
			t1.Fail()
		}
	})

	t.Run("Error tracks not exist", func(t2 *testing.T) {
		count := uint64(3)
		userID := uint64(1)
		trackRepo.EXPECT().FetchFavourites(userID, count).Return(nil, nil)
		usecase := configTrackUsecase(trackRepo)

		expTracks, err := exe(usecase, userID, count)

		if gomock.Not([]*Track{}).Matches(expTracks) || err != nil {
			fmt.Println("Error. Expected: nil, Got: ", err)
			fmt.Println("expTracks. Expected:", []Tracks{} , "Got:", expTracks)
			t2.Fail()
		}
	})

	t.Run("Error fetching", func(t3 *testing.T) {
		count := uint64(3)
		userID := uint64(1)
		trackRepo.EXPECT().FetchFavourites(userID, count).Return(nil, fmt.Errorf("some error"))
		usecase := configTrackUsecase(trackRepo)

		expTracks, err := exe(usecase, userID, count)

		if expTracks != nil || err == nil {
			fmt.Println("Error. Expected: not nil, Got: ", err)
			fmt.Println("expTracks. Expected: nil, Got: ", expTracks)
			t3.Fail()
		}
	})
}
