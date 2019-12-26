package search

import "2019_2_Covenant/internal/models"

type Usecase interface {
	Search(text string, count uint64, authID uint64) ([]*models.Track, []*models.Album, []*models.Artist, error)
}
