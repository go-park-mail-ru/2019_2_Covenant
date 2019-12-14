package usecase

import (
	"2019_2_Covenant/internal/album"
	"2019_2_Covenant/internal/artist"
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/search"
	"2019_2_Covenant/internal/track"
	. "2019_2_Covenant/tools/vars"
)

type SearchUsecase struct {
	trackRepo  track.Repository
	albumRepo  album.Repository
	artistRepo artist.Repository
}

func NewSearchUsecase(tR track.Repository, alR album.Repository, arR artist.Repository) search.Usecase {
	return &SearchUsecase{
		trackRepo:  tR,
		albumRepo:  alR,
		artistRepo: arR,
	}
}

func (su *SearchUsecase) Search(text string, count uint64) ([]*models.Track, []*models.Album, []*models.Artist, error) {
	tracks, err := su.trackRepo.FindLike(text, count)
	if err != nil {
		return nil, nil, nil, err
	}

	albums, err := su.albumRepo.FindLike(text, count)
	if err != nil {
		return nil, nil, nil, err
	}

	artists, err := su.artistRepo.FindLike(text, count)
	if err != nil {
		return nil, nil, nil, err
	}

	if tracks == nil && albums == nil && artists == nil {
		return nil, nil, nil, ErrNotFound
	}

	if tracks == nil { tracks = []*models.Track{} }
	if albums == nil { albums = []*models.Album{} }
	if artists == nil { artists = []*models.Artist{} }

	return tracks, albums, artists, nil
}
