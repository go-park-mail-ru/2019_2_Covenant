package storage

type Storage interface {
	Open() error
	Close()
}
