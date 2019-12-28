package storage

import (
	"2019_2_Covenant/internal/album"
	_albumRepo "2019_2_Covenant/internal/album/repository"
	"2019_2_Covenant/internal/artist"
	_artistRepo "2019_2_Covenant/internal/artist/repository"
	"2019_2_Covenant/internal/collections"
	_collectionRepo "2019_2_Covenant/internal/collections/repository"
	"2019_2_Covenant/internal/likes"
	_likesRepo "2019_2_Covenant/internal/likes/repository"
	"2019_2_Covenant/internal/playlist"
	_playlistRepo "2019_2_Covenant/internal/playlist/repository"
	"2019_2_Covenant/internal/session"
	_sessRepo "2019_2_Covenant/internal/session/repository"
	"2019_2_Covenant/internal/subscriptions"
	_subscriptionRepo "2019_2_Covenant/internal/subscriptions/repository"
	"2019_2_Covenant/internal/track"
	_trackRepo "2019_2_Covenant/internal/track/repository"
	"2019_2_Covenant/internal/user"
	_userRepo "2019_2_Covenant/internal/user/repository"
	"database/sql"
	_ "github.com/lib/pq"
)

type BaseStorage struct {
	config *Config
	db     *sql.DB
}

type PGStorage struct {
	BaseStorage
	userRepo         user.Repository
	sessRepo         session.Repository
	trackRepo        track.Repository
	playlistRepo     playlist.Repository
	albumRepo        album.Repository
	artistRepo       artist.Repository
	subscriptionRepo subscriptions.Repository
	likesRepo        likes.Repository
	collectionRepo   collections.Repository
}

func NewPGStorage(conf *Config) Storage {
	return &PGStorage{
		BaseStorage: BaseStorage{
			config: conf,
		},
	}
}

func (s *PGStorage) Open() error {
	db, err := sql.Open("postgres", s.config.GetURL())

	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	s.db = db

	return nil
}

func (s *PGStorage) Close() {
	s.db.Close()
}

func (s *PGStorage) User() user.Repository {
	if s.userRepo != nil {
		return s.userRepo
	}

	s.userRepo = _userRepo.NewUserRepository(s.db)

	return s.userRepo
}

func (s *PGStorage) Session() session.Repository {
	if s.sessRepo != nil {
		return s.sessRepo
	}

	s.sessRepo = _sessRepo.NewSessionRepository(s.db)

	return s.sessRepo
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

func (s *PGStorage) Collection() collections.Repository {
	if s.collectionRepo != nil {
		return s.collectionRepo
	}

	s.collectionRepo = _collectionRepo.NewCollectionRepository(s.db)

	return s.collectionRepo
}
