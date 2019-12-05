package subscriptions

type Repository interface {
	Subscribe(userID uint64, subscribedID uint64) error
	Unsubscribe(userID uint64, subscribedID uint64) error
}
