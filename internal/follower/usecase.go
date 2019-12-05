package follower

import "2019_2_Covenant/internal/models"

type Usecase interface {
	Subscribe(userID uint64, followerId uint64) error
	Unsubscribe(userID uint64, followerId uint64) error
	GetProfile(userID uint64) (*models.User, error)
}

