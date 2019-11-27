package artist

import "2019_2_Covenant/internal/models"

type Usecase interface {
	FindLike(name string, count uint64) ([]*models.Artist, error)
}
