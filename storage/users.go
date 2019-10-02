package storage

import (
	"fmt"
	"sync"
)

type User struct {
	ID uint `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	Avatar string `json:"avatar"`
}

type UserStore struct {
	users  []*User
	mu     sync.RWMutex
	nextID uint
}

func NewUserStore() *UserStore {
	return &UserStore{
		mu:    sync.RWMutex{},
		users: []*User{},
	}
}

func (us *UserStore) AddUser(newUser *User) (uint, error) {
	us.mu.Lock()
	defer us.mu.Unlock()

	newUser.ID = us.nextID
	us.nextID++

	us.users = append(us.users, newUser)

	return newUser.ID, nil
}

func (us *UserStore) GetUserByID(userID uint) (*User, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()

	for _, user := range us.users {
		if user.ID == userID {
			return user, nil
		}
	}

	err := fmt.Errorf("user is not exist")
	return nil, err
}

func (us *UserStore) GetUserByUsername(username string) (*User, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()

	for _, user := range us.users {
		if user.Username == username {
			return user, nil
		}
	}

	err := fmt.Errorf("user is not exist")
	return nil, err
}

func (us *UserStore) CheckUser(email string, password string) (uint, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()

	for _, user := range us.users {
		if user.Email == email && user.Password == password {
			return user.ID, nil
		}
	}

	err := fmt.Errorf("email and password are mismatched")
	return 0, err
}