package usecase

import (
	"2019_2_Covenant/internal/artist"
	"2019_2_Covenant/internal/models"
)

type ArtistUsecase struct {
	artistRepo artist.Repository
}

func NewAlbumUsecase(repo artist.Repository) artist.Usecase {
	return &ArtistUsecase{
		artistRepo: repo,
	}
}

func (aUC *ArtistUsecase) FindLike(name string, count uint64) ([]*models.Artist, error) {
	albums, err := aUC.artistRepo.FindLike(name, count)

	if err != nil {
		return nil, err
	}

	if albums == nil {
		albums = []*models.Artist{}
	}

	return albums, nil
}
