package user

import (
	. "../models"
)

/*
 *	Repository interface represents the user's repository contract
 */

type Repository interface {
	Fetch(number int64) ([]*User, error)
	GetById(id int64) (*User, error)
	GetByEmail(email string) (*User, error)
	Store (user *User) error
}
