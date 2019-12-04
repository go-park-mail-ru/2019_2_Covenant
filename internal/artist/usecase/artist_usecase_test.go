package usecase

import (
	"2019_2_Covenant/internal/artist"
	mock "2019_2_Covenant/internal/artist/mocks"
	. "2019_2_Covenant/internal/models"
	"fmt"
	"github.com/golang/mock/gomock"
	"testing"
)

//go:generate mockgen -source=../repository.go -destination=../mocks/mock_repository.go -package=mock

type Artists struct {
	Artist []*Artist
}

var artists = Artists{
	Artist: []*Artist{
		{ID: 1, Name: "News of the World"},
		{ID: 2, Name: "WHEN WE ALL FALL ASLEEP, WHERE DO WE GO?"},
		{ID: 3, Name: "Love at First Sting"},
	},
}

func configArtistUsecase(artistRepo *mock.MockRepository) artist.Usecase {
	return NewArtistUsecase(artistRepo)
}

func TestArtistUsecase_FindLike(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	artistRepo := mock.NewMockRepository(ctrl)

	exe := func(usecase artist.Usecase, findStr string, count uint64) ([]*Artist, error) {
		return usecase.FindLike(findStr, count)
	}

	t.Run("Test OK", func(t1 *testing.T) {
		count := uint64(3)
		findStr := "f"

		artistRepo.EXPECT().FindLike(findStr, count).Return(artists.Artist, nil)
		usecase := configArtistUsecase(artistRepo)

		expArtists, err := exe(usecase, findStr, count)

		if gomock.Not(artists.Artist).Matches(expArtists) || err != nil {
			fmt.Println("Error. Expected: nil, Got:", err)
			fmt.Println("expArtists. Expected:", artists.Artist, "Got:", expArtists)
			t1.Fail()
		}
	})

	t.Run("Error artists not exist", func(t2 *testing.T) {
		count := uint64(3)
		findStr := "f"

		artistRepo.EXPECT().FindLike(findStr, count).Return(nil, nil)
		usecase := configArtistUsecase(artistRepo)

		expArtists, err := exe(usecase, findStr, count)

		if gomock.Not([]*Artist{}).Matches(expArtists) || err != nil {
			fmt.Println("Error. Expected: nil, Got:", err)
			fmt.Println("expArtists. Expected:", []Album{}, "Got:", expArtists)
			t2.Fail()
		}
	})

	t.Run("Error fetching", func(t3 *testing.T) {
		count := uint64(3)
		findStr := "f"

		artistRepo.EXPECT().FindLike(findStr, count).Return(nil, fmt.Errorf("some error"))
		usecase := configArtistUsecase(artistRepo)

		expArtists, err := exe(usecase, findStr, count)

		if expArtists != nil || err == nil {
			fmt.Println("Error. Expected: not nil, Got:", err)
			fmt.Println("expArtists. Expected: nil, Got:", expArtists)
			t3.Fail()
		}
	})
}

