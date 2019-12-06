package subscriptions

type Repository interface {
	Subscribe(userID uint64, subscriptionID uint64) error
	Unsubscribe(userID uint64, subscriptionID uint64) error
}
