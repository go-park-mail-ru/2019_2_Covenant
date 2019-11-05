package track

import "2019_2_Covenant/internal/models"

type Repository interface {
	Fetch(count uint64) ([]*models.Track, error)
}
