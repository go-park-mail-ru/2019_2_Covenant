package usecase

import (
	"2019_2_Covenant/pkg/album"
	files "2019_2_Covenant/pkg/file_processor"
	"2019_2_Covenant/pkg/models"
	"context"
	"io"
)

type AlbumUsecase struct {
	albumRepo album.Repository
	filesRepo files.Repository
}

func NewAlbumUsecase(repo album.Repository, fr files.Repository) album.Usecase {
	return &AlbumUsecase{
		albumRepo: repo,
		filesRepo: fr,
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

	for _, a := range albums {
		a.Year = a.Year[:4]
	}

	return albums, nil
}

func (aUC *AlbumUsecase) DeleteByID(id uint64) error {
	if err := aUC.albumRepo.DeleteByID(id); err != nil {
		return err
	}

	return nil
}

func (aUC *AlbumUsecase) UpdateByID(albumID uint64, artistID uint64, name string, year string) error {
	if err := aUC.albumRepo.UpdateByID(albumID, artistID, name, year); err != nil {
		return err
	}

	return nil
}

func (aUC *AlbumUsecase) Fetch(count uint64, offset uint64) ([]*models.Album, uint64, error) {
	albums, total, err := aUC.albumRepo.Fetch(count, offset)

	if err != nil {
		return nil, total, err
	}

	if albums == nil {
		albums = []*models.Album{}
	}

	for _, a := range albums {
		a.Year = a.Year[:4]
	}

	return albums, total, nil
}

func (aUC *AlbumUsecase) GetByID(id uint64) (*models.Album, uint64, error) {
	a, amountOfTracks, err := aUC.albumRepo.GetByID(id)

	if err != nil {
		return nil, amountOfTracks, err
	}

	a.Year = a.Year[:4]

	return a, amountOfTracks, nil
}

func (aUC *AlbumUsecase) AddTrack(ctx context.Context, albumID uint64, name string, track io.Reader) error {
	id, err := aUC.albumRepo.AddTrack(albumID, name)
	if err != nil {
		return err
	}

	return aUC.filesRepo.ProcessTrack(ctx, track, id)
}

func (aUC *AlbumUsecase) GetTracksFrom(albumID uint64, authID uint64) ([]*models.Track, error) {
	tracks, err := aUC.albumRepo.GetTracksFrom(albumID, authID)

	if err != nil {
		return nil, err
	}

	if tracks == nil {
		tracks = []*models.Track{}
	}

	return tracks, nil
}

func (aUC *AlbumUsecase) UpdatePhoto(ctx context.Context, albumID uint64, photo io.Reader) error {
	return aUC.filesRepo.ProcessAlbumPhoto(ctx, photo, albumID)
}
