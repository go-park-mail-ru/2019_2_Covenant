package repository

import (
	. "2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/session"
	"2019_2_Covenant/internal/vars"
	"sync"
	"time"
)

type SessionStorage struct {
	sessions []*Session
	mu       sync.RWMutex
	nextID   uint64
}

func NewSessionStorage() session.Repository {
	return &SessionStorage{
		mu:       sync.RWMutex{},
		sessions: []*Session{},
	}
}

func (ss *SessionStorage) Get(value string) (*Session, error) {
	item := &Session{}
	var isFound bool

	for i := 0; i < len(ss.sessions) && !isFound; i++ {
		if ss.sessions[i].Data == value {
			item = ss.sessions[i]
			isFound = true
		}
	}

	if !isFound {
		return nil, vars.ErrNotFound
	}

	timeNow := time.Now()
	diffTime := item.Expires.Sub(timeNow)

	if diffTime <= 0 {
		err := ss.DeleteByID(item.ID)

		if err != nil {
			return nil, err
		}

		return nil, vars.ErrExpired
	}

	return item, nil
}

func (ss *SessionStorage) Store(newSession *Session) error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	newSession.ID = ss.nextID
	ss.nextID++

	ss.sessions = append(ss.sessions, newSession)

	return nil
}

func (ss *SessionStorage) DeleteByID(id uint64) error {
	remove := func(id uint64) ([]*Session, error) {
		if len(ss.sessions) == 0 {
			return nil, vars.ErrNotFound
		}
		ss.sessions[len(ss.sessions)-1], ss.sessions[id] = ss.sessions[id], ss.sessions[len(ss.sessions)-1]
		return ss.sessions[:len(ss.sessions)-1], nil
	}

	var err error
	ss.sessions, err = remove(id)

	if err != nil {
		return err
	}

	return nil
}
