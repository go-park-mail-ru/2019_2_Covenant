package usecase

import (
	"2019_2_Covenant/internal/follower"
)

type FollowerUsecase struct {
	followerRepo follower.Repository
}

func NewFollowerUsecase(repo follower.Repository) follower.Usecase {
	return &FollowerUsecase{
		followerRepo: repo,
	}
}

func (fUc *FollowerUsecase) Subscribe(userID uint64, followerId uint64) error {
	return fUc.followerRepo.Subscribe(userID, followerId)
}

func (fUc *FollowerUsecase) Unsubscribe(userID uint64, followerId uint64) error {
	return fUc.followerRepo.Unsubscribe(userID, followerId)
}

