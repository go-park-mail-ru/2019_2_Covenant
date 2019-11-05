package track

import "2019_2_Covenant/internal/models"

type Usecase interface {
	Fetch(count uint64) ([]*models.Track, error)
}
