package usecase

import (
	"2019_2_Covenant/internal/album"
	mock "2019_2_Covenant/internal/album/mocks"
	. "2019_2_Covenant/internal/models"
	"fmt"
	"github.com/golang/mock/gomock"
	"testing"
)


//go:generate mockgen -source=../repository.go -destination=../mocks/mock_repository.go -package=mock

type Albums struct {
	Album []*Album
}

var albums = Albums{
	Album: []*Album{
		{ID: 1, Name: "News of the World"},
		{ID: 2, Name: "WHEN WE ALL FALL ASLEEP, WHERE DO WE GO?"},
		{ID: 3, Name: "Love at First Sting"},
	},
}

func configAlbumUsecase(albumRepo *mock.MockRepository) album.Usecase {
	return NewAlbumUsecase(albumRepo)
}

func TestAlbumUsecase_FindLike(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	albumRepo := mock.NewMockRepository(ctrl)

	exe := func(usecase album.Usecase, findStr string, count uint64) ([]*Album, error) {
		return usecase.FindLike(findStr, count)
	}

	t.Run("Test OK", func(t1 *testing.T) {
		count := uint64(3)
		findStr := "f"

		albumRepo.EXPECT().FindLike(findStr, count).Return(albums.Album, nil)
		usecase := configAlbumUsecase(albumRepo)

		expAlbums, err := exe(usecase, findStr, count)

		if gomock.Not(albums.Album).Matches(expAlbums) || err != nil {
			fmt.Println("Error. Expected: nil, Got:", err)
			fmt.Println("expAlbums. Expected:", albums.Album, "Got:", expAlbums)
			t1.Fail()
		}
	})

	t.Run("Error albums not exist", func(t2 *testing.T) {
		count := uint64(3)
		findStr := "f"

		albumRepo.EXPECT().FindLike(findStr, count).Return(nil, nil)
		usecase := configAlbumUsecase(albumRepo)

		expAlbums, err := exe(usecase, findStr, count)

		if gomock.Not([]*Album{}).Matches(expAlbums) || err != nil {
			fmt.Println("Error. Expected: nil, Got:", err)
			fmt.Println("expAlbums. Expected:", []Album{}, "Got:", expAlbums)
			t2.Fail()
		}
	})

	t.Run("Error fetching", func(t3 *testing.T) {
		count := uint64(3)
		findStr := "f"

		albumRepo.EXPECT().FindLike(findStr, count).Return(nil, fmt.Errorf("some error"))
		usecase := configAlbumUsecase(albumRepo)

		expAlbums, err := exe(usecase, findStr, count)

		if expAlbums != nil || err == nil {
			fmt.Println("Error. Expected: not nil, Got:", err)
			fmt.Println("expAlbums. Expected: nil, Got:", expAlbums)
			t3.Fail()
		}
	})
}
