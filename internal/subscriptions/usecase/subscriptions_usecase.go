package usecase

import (
	"2019_2_Covenant/internal/subscriptions"
	. "2019_2_Covenant/tools/vars"
)

type SubscriptionUsecase struct {
	subscriptionRepo subscriptions.Repository
}

func NewSubscriptionUsecase(repo subscriptions.Repository) subscriptions.Usecase {
	return &SubscriptionUsecase{
		subscriptionRepo: repo,
	}
}

func (fUc *SubscriptionUsecase) Subscribe(userID uint64, subscriptionID uint64) error {
	err := fUc.subscriptionRepo.Subscribe(userID, subscriptionID)

	if err == ErrAlreadyExist {
		return err
	}

	if err != nil {
		return ErrInternalServerError
	}

	return nil
}

func (fUc *SubscriptionUsecase) Unsubscribe(userID uint64, subscriptionID uint64) error {
	err := fUc.subscriptionRepo.Unsubscribe(userID, subscriptionID)

	if err == ErrNotFound {
		return err
	}

	if err != nil {
		return ErrInternalServerError
	}

	return nil
}

