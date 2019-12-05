package follower

import "2019_2_Covenant/internal/models"

type Repository interface {
	Subscribe(userID uint64, followerId uint64) error
	Unsubscribe(userID uint64, followerId uint64) error
	GetProfile(usrID uint64) (*models.User, error)
}
