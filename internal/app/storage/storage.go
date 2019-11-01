package storage

import (
	"2019_2_Covenant/internal/session"
	_sessRepo "2019_2_Covenant/internal/session/repository"
	"2019_2_Covenant/internal/user"
	_userRepo "2019_2_Covenant/internal/user/repository"
	"database/sql"
	_ "github.com/lib/pq"
)

type BaseStorage struct {
	config *Config
	db *sql.DB
}

type PGStorage struct {
	BaseStorage
	userRepo user.Repository
	sessRepo session.Repository
}

func NewPGStorage(conf *Config) Storage {
	return &PGStorage{
		BaseStorage: BaseStorage{
			config: conf,
		},
	}
}

func (s *PGStorage) Open() error {
	db, err := sql.Open("postgres", s.config.DBUrl)

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

	s.sessRepo = _sessRepo.NewSessionStorage()

	return s.sessRepo
}
