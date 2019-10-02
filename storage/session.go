package storage

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Session struct {
	UserID uint
	Expires time.Time
}

type SessionsStore struct {
	sessions  map[string]*Session
	mu     sync.RWMutex
}

func generateSessionID(length uint) (sessionID string) {
	rand.Seed(time.Now().UnixNano())
	chars := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

	for i := uint(0); i < length; i++ {
		sessionID += string(chars[rand.Intn(len(chars))])
	}
	return
}

func (s *SessionsStore) Set(userID uint) (sessionID string, session *Session) {
	loc, _ := time.LoadLocation("Europe/Moscow")

	expires := time.Now().In(loc).Add(24 * time.Hour)

	session = &Session{
		UserID: userID,
		Expires: expires,
	}

	sessionID = generateSessionID(10)

	s.mu.Lock()
	s.sessions[sessionID] = session
	s.mu.Unlock()

	return
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