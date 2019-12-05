package usecase

import (
	"2019_2_Covenant/internal/subscriptions"
)

type SubscriptionUsecase struct {
	subscriptionRepo subscriptions.Repository
}

func NewSubscriptionUsecase(repo subscriptions.Repository) subscriptions.Usecase {
	return &SubscriptionUsecase{
		subscriptionRepo: repo,
	}
}

func (fUc *SubscriptionUsecase) Subscribe(userID uint64, subscribedID uint64) error {
	return fUc.subscriptionRepo.Subscribe(userID, subscribedID)
}

func (fUc *SubscriptionUsecase) Unsubscribe(userID uint64, subscribedID uint64) error {
	return fUc.subscriptionRepo.Unsubscribe(userID, subscribedID)
}

