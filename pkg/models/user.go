package models

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"golang.org/x/crypto/pbkdf2"
)

type User struct {
	ID            uint64 `json:"id"`
	Nickname      string `json:"nickname"`
	Email         string `json:"email"`
	PlainPassword string `json:"-"`
	Password      string `json:"-"`
	Avatar        string `json:"avatar"`
	Role          int8   `json:"role"`   // 0 - user; 1 - admin;
	Access        int8   `json:"access"` // 0 - public; 1 - private;
	Subscription  *bool  `json:"subscription,omitempty"`
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

	return base64.StdEncoding.EncodeToString(result), nil
}

func (u *User) Verify(plainPassword string) bool {
	encryptedPassword, _ := base64.StdEncoding.DecodeString(u.Password)
	salt := encryptedPassword[0:8]
	dk := pbkdf2.Key([]byte(plainPassword), salt, 4096, 32, sha1.New)

	return bytes.Equal(dk, encryptedPassword[8:])
}
