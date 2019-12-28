package subscriptions

type Usecase interface {
	Subscribe(userID uint64, subscriptionID uint64) error
	Unsubscribe(userID uint64, subscriptionID uint64) error
}
