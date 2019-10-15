package repository

import (
	. "2019_2_Covenant/pkg/models"
	"2019_2_Covenant/pkg/user"
	"fmt"
	"sync"
)

type UserStorage struct {
	users  []*User
	mu     sync.RWMutex
	nextID uint64
}

func NewUserStorage() user.Repository {
	return &UserStorage{
		mu:    sync.RWMutex{},
		users: []*User{},
	}
}

func (us *UserStorage) Store(newUser *User) error {
	us.mu.Lock()
	defer us.mu.Unlock()

	newUser.ID = us.nextID
	us.nextID++

	newUser.Username = "User" + fmt.Sprint(newUser.ID)
	newUser.Avatar = "img/user_profile.png"

	us.users = append(us.users, newUser)

	return nil
}

func (us *UserStorage) GetByEmail(email string) (*User, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()

	for _, user := range us.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, ErrNotFound
}

func (us *UserStorage) GetByID(userID uint64) (*User, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()

	for _, user := range us.users {
		if user.ID == userID {
			return user, nil
		}
	}

	return nil, ErrNotFound
}

func (us *UserStorage) FetchAll() ([]*User, error) {
	return us.users, nil
}
