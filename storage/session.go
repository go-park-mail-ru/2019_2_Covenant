package storage

import (
	"sync"
	"time"
)

//cp -rf ~/Documents/covenant-backend/2019_2_Covenant ~/go/src

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Session struct {
	ID 		  uint
	UserID    uint
	Expires   time.Time
	Data      string
}

type SessionsStore struct {
	sessions  map[string]*Session
	mu     	  sync.RWMutex
	nextID    uint
}

func NewSessionStore() *SessionsStore {
	return &SessionsStore{
		mu:    sync.RWMutex{},
		sessions: map[string]*Session{},
	}
}

func generateSessionID(length uint) (sessionID string) {
	rand.Seed(time.Now().UnixNano())
	chars := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

	for i := uint(0); i < length; i++ {
		sessionID += string(chars[rand.Intn(len(chars))])
	}

	return sessionID
}

func (s *SessionsStore) Set(newSession *Session) (uint, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	newSession.ID = s.nextID
	s.nextID++

	s.sessions[newSession.Data] = newSession

	return newSession.ID, nil
}

func (s *SessionsStore) Get(sessionID string) (session *Session, err error) {
	loc, _ := time.LoadLocation("Europe/Moscow")

	s.mu.RLock()
	session, ok := s.sessions[sessionID]
	s.mu.RUnlock()

	if !ok {
		err = fmt.Errorf("session does not exist")
		return
	}

	timeNow := time.Now().In(loc)
	diffTime := session.Expires.Sub(timeNow)

	if diffTime < 0 {
		s.Delete(sessionID)
		err = fmt.Errorf("session expired")
		return
	}

	return
}

func (s *SessionsStore) Delete(sessionID string) {
	s.mu.Lock()
	delete(s.sessions, sessionID)
	s.mu.Unlock()
}
