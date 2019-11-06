package models

import (
	"2019_2_Covenant/internal/vars"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Session struct {
	ID      uint64
	UserID  uint64
	Expires time.Time
	Data    string
}

type CSRFTokenManager struct {
	Secret []byte
}

func NewCSRFTokenManager(secret string) *CSRFTokenManager {
	return &CSRFTokenManager{
		Secret: []byte(secret),
	}
}

func (tk *CSRFTokenManager) Create(usr *User, sess *Session, expires time.Time) (string, error) {
	h := hmac.New(sha1.New, tk.Secret)
	data := fmt.Sprintf("%s:%s:%d", usr.ID, sess.Data, expires.Unix())
	h.Write([]byte(data))

	token := hex.EncodeToString(h.Sum(nil)) + ":" + strconv.FormatInt(expires.Unix(), 10)

	return token, nil
}

func (tk *CSRFTokenManager) Verify(usr *User, sess *Session, token string) (bool, error) {
	tokenData := strings.Split(token, ":")

	if len(tokenData) != 2 {
		return false, vars.ErrExpired
	}

	tokenExp, err := strconv.ParseInt(tokenData[1], 10, 64)

	if err != nil {
		return false, vars.ErrInternalServerError
	}

	if tokenExp < time.Now().Unix() {
		return false, vars.ErrExpired
	}

	h := hmac.New(sha1.New, tk.Secret)
	data := fmt.Sprintf("%s:%s:%d", usr.ID, sess.Data, tokenExp)
	h.Write([]byte(data))

	expected := h.Sum(nil)
	got, err := hex.DecodeString(tokenData[0])

	if err != nil {
		return false, vars.ErrInternalServerError
	}

	return hmac.Equal(expected, got), nil
}
