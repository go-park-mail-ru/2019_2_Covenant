package usecase

import (
	"2019_2_Covenant/internal/album"
	"2019_2_Covenant/internal/models"
)

type AlbumUsecase struct {
	albumRepo album.Repository
}

func NewAlbumUsecase(repo album.Repository) album.Usecase {
	return &AlbumUsecase{
		albumRepo: repo,
	}
}

func (aUC *AlbumUsecase) FindLike(name string, count uint64) ([]*models.Album, error) {
	albums, err := aUC.albumRepo.FindLike(name, count)

	if err != nil {
		return nil, err
	}

	if albums == nil {
		albums = []*models.Album{}
	}

	return albums, nil
}
