package repository

import (
	. "2019_2_Covenant/pkg/models"
	"2019_2_Covenant/pkg/user"
	"2019_2_Covenant/pkg/vars"
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

	for _, usr := range us.users {
		if usr.Email == email {
			return usr, nil
		}
	}

	return nil, vars.ErrNotFound
}

func (us *UserStorage) GetByID(usrID uint64) (*User, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()

	for _, usr := range us.users {
		if usr.ID == usrID {
			return usr, nil
		}
	}

	return nil, vars.ErrNotFound
}

func (us *UserStorage) FetchAll() ([]*User, error) {
	return us.users, nil
}
