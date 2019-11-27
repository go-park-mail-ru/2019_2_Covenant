package models

import (
	"crypto/rand"
	"crypto/sha1"

	"golang.org/x/crypto/pbkdf2"
)

type User struct {
	ID            uint64 `json:"-"`
	Nickname      string `json:"nickname"`
	Email         string `json:"email"`
	PlainPassword string `json:"-"`
	Password      string `json:"-"`
	Avatar        string `json:"avatar"`
	Role          int8   `json:"role"`   // 0 - user; 1 - admin;
	Access        int8   `json:"access"` // 0 - public; 1 - private;
}

func NewUser(email string, nickname string, plainPassword string) *User {
	return &User{
		Nickname:      nickname,
		Email:         email,
		PlainPassword: plainPassword,
	}
}

func (u *User) BeforeStore() error {
	if len(u.PlainPassword) > 0 {
		pass, err := EncryptPassword(u.PlainPassword)

		if err != nil {
			return err
		}

		u.Password = pass
	}

	return nil
}

func EncryptPassword(plainPassword string) (string, error) {
	salt := make([]byte, 8)

	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	dk := pbkdf2.Key([]byte(plainPassword), salt, 4096, 32, sha1.New)

	result := append(salt, dk...)

	return string(result), nil
}

func (u *User) Verify(plainPassword string) bool {
	salt := u.Password[0:8]
	dk := pbkdf2.Key([]byte(plainPassword), []byte(salt), 4096, 32, sha1.New)
	return string(dk) == u.Password[8:]
}
