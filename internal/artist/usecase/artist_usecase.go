package usecase

import (
	"2019_2_Covenant/internal/artist"
	"2019_2_Covenant/internal/models"
)

type ArtistUsecase struct {
	artistRepo artist.Repository
}

func NewArtistUsecase(repo artist.Repository) artist.Usecase {
	return &ArtistUsecase{
		artistRepo: repo,
	}
}

func (aUC *ArtistUsecase) FindLike(name string, count uint64) ([]*models.Artist, error) {
	artists, err := aUC.artistRepo.FindLike(name, count)

	if err != nil {
		return nil, err
	}

	if artists == nil {
		artists = []*models.Artist{}
	}

	return artists, nil
}

func (aUC *ArtistUsecase) Store(artist *models.Artist) error {
	if err := aUC.artistRepo.Store(artist); err != nil {
		return err
	}

	return nil
}

func (aUC *ArtistUsecase) DeleteByID(id uint64) error {
	if err := aUC.artistRepo.DeleteByID(id); err != nil {
		return err
	}

	return nil
}

func (aUC *ArtistUsecase) UpdateByID(id uint64, name string) error {
	if err := aUC.artistRepo.UpdateByID(id, name); err != nil {
		return err
	}

	return nil
}

func (aUC *ArtistUsecase) Fetch(count uint64, offset uint64) ([]*models.Artist, uint64, error) {
	artists, total, err := aUC.artistRepo.Fetch(count, offset)

	if err != nil {
		return nil, total, err
	}

	if artists == nil {
		artists = []*models.Artist{}
	}

	return artists, total, nil
}

func (aUC *ArtistUsecase) CreateAlbum(album *models.Album) error {
	if err := aUC.artistRepo.CreateAlbum(album); err != nil {
		return err
	}

	return nil
}

func (aUC *ArtistUsecase) GetByID(id uint64) (*models.Artist, uint64, error) {
	a, amountOfAlbums, err := aUC.artistRepo.GetByID(id)

	if err != nil {
		return nil, amountOfAlbums, err
	}

	return a, amountOfAlbums, nil
}
