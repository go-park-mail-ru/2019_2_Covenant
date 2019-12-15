package usecase

import (
	"2019_2_Covenant/internal/artist"
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/tools/time_parser"
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

func (aUC *ArtistUsecase) UpdatePhoto(artistID uint64, path string) error {
	if err := aUC.artistRepo.UpdatePhoto(artistID, path); err != nil {
		return err
	}

	return nil
}

func (aUC *ArtistUsecase) GetArtistAlbums(artistID uint64, count uint64, offset uint64) ([]*models.Album, uint64, error) {
	albums, total, err := aUC.artistRepo.GetArtistAlbums(artistID, count, offset)

	if err != nil {
		return nil, total, err
	}

	if albums == nil {
		albums = []*models.Album{}
	}

	for _, a := range albums { a.Year = a.Year[:4] }

	return albums, total, nil
}

func (aUC *ArtistUsecase) GetTracks(artistID uint64, count uint64, offset uint64) ([]*models.Track, uint64, error) {
	tracks, total, err := aUC.artistRepo.GetTracks(artistID, count, offset)

	if err != nil {
		return nil, total, err
	}

	if tracks == nil {
		tracks = []*models.Track{}
	}

	for _, item := range tracks { item.Duration = time_parser.GetDuration(item.Duration) }

	return tracks, total, nil
}
