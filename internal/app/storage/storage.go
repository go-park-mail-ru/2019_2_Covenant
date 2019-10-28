package storage

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Storage struct {
	config *Config
	db *sql.DB
}

func NewStorage(conf *Config) *Storage {
	return &Storage{
		config: conf,
	}
}

func (s *Storage) Open() error {
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

func (s *Storage) Close() {
	s.db.Close()
}
