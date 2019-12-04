package follower

type Usecase interface {
	Subscribe(userID uint64, followerId uint64) error
	Unsubscribe(userID uint64, followerId uint64) error
}

