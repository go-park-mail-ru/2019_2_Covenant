package storage

import (
	"2019_2_Covenant/internal/album"
	"2019_2_Covenant/internal/artist"
	"2019_2_Covenant/internal/playlist"
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
	Playlist() playlist.Repository
	Album() album.Repository
	Artist() artist.Repository
}
