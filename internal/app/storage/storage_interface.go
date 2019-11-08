package storage

import (
	"2019_2_Covenant/internal/session"
	"2019_2_Covenant/internal/track"
	"2019_2_Covenant/internal/user"
)

type Storage interface {
	Open() error
	Close()
	User() user.Repository
	Session() session.Repository
	Track() track.Repository
}
