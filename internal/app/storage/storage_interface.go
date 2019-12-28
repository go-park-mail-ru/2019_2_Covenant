package storage

import (
	"2019_2_Covenant/internal/album"
	"2019_2_Covenant/internal/artist"
	"2019_2_Covenant/internal/collections"
	"2019_2_Covenant/internal/likes"
	"2019_2_Covenant/internal/subscriptions"
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
	Subscription() subscriptions.Repository
	Like() likes.Repository
	Collection() collections.Repository
}
