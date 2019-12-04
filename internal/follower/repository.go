package follower

type Repository interface {
	Subscribe(userID uint64, followerId uint64) error
	Unsubscribe(userID uint64, followerId uint64) error
}
