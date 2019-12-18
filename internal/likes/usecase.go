package likes

type Usecase interface {
	Like(userID uint64, trackID uint64) error
	Unlike(userID uint64, trackID uint64) error
}
