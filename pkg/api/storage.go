package api

import (
	"2019_2_Covenant/pkg/album"
	_albumRepo "2019_2_Covenant/pkg/album/repository"
	"2019_2_Covenant/pkg/artist"
	_artistRepo "2019_2_Covenant/pkg/artist/repository"
	"2019_2_Covenant/pkg/likes"
	_likesRepo "2019_2_Covenant/pkg/likes/repository"
	"2019_2_Covenant/pkg/playlist"
	_playlistRepo "2019_2_Covenant/pkg/playlist/repository"
	"2019_2_Covenant/pkg/subscriptions"
	_subscriptionRepo "2019_2_Covenant/pkg/subscriptions/repository"
	"2019_2_Covenant/pkg/track"
	_trackRepo "2019_2_Covenant/pkg/track/repository"
	"database/sql"
)

type PGStorage struct {
	db               *sql.DB
	trackRepo        track.Repository
	playlistRepo     playlist.Repository
	albumRepo        album.Repository
	artistRepo       artist.Repository
	subscriptionRepo subscriptions.Repository
	likesRepo        likes.Repository
}

func NewPGStorage(db *sql.DB) *PGStorage {
	return &PGStorage{
		db: db,
	}
}

func (s *PGStorage) Track() track.Repository {
	if s.trackRepo != nil {
		return s.trackRepo
	}

	s.trackRepo = _trackRepo.NewTrackRepository(s.db)

	return s.trackRepo
}

func (s *PGStorage) Playlist() playlist.Repository {
	if s.playlistRepo != nil {
		return s.playlistRepo
	}

	s.playlistRepo = _playlistRepo.NewPlaylistRepository(s.db)

	return s.playlistRepo
}

func (s *PGStorage) Album() album.Repository {
	if s.albumRepo != nil {
		return s.albumRepo
	}

	s.albumRepo = _albumRepo.NewAlbumRepository(s.db)

	return s.albumRepo
}

func (s *PGStorage) Artist() artist.Repository {
	if s.artistRepo != nil {
		return s.artistRepo
	}

	s.artistRepo = _artistRepo.NewArtistRepository(s.db)

	return s.artistRepo
}

func (s *PGStorage) Subscription() subscriptions.Repository {
	if s.subscriptionRepo != nil {
		return s.subscriptionRepo
	}

	s.subscriptionRepo = _subscriptionRepo.NewSubscriptionRepository(s.db)

	return s.subscriptionRepo
}

func (s *PGStorage) Like() likes.Repository {
	if s.likesRepo != nil {
		return s.likesRepo
	}

	s.likesRepo = _likesRepo.NewLikesRepository(s.db)

	return s.likesRepo
}
