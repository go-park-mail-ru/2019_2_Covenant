package models

import (
	"2019_2_Covenant/internal/vars"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"net/http"
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

func NewSession(userID uint64) (*Session, *http.Cookie) {
	cookie := &http.Cookie{
		Name:    "Covenant",
		Value:   uuid.New().String(),
		Expires: time.Now().Add(24 * time.Hour),
	}

	return &Session{
		UserID:  userID,
		Data:    cookie.Value,
		Expires: cookie.Expires,
	}, cookie
}

type CSRFTokenManager struct {
	Secret []byte
}

func NewCSRFTokenManager(secret string) *CSRFTokenManager {
	return &CSRFTokenManager{
		Secret: []byte(secret),
	}
}

func (tk *CSRFTokenManager) Create(user_id uint64, cookie string, expires time.Time) (string, error) {
	h := hmac.New(sha1.New, tk.Secret)
	data := fmt.Sprintf("%d:%s:%d", user_id, cookie, expires.Unix())

	if _, err := h.Write([]byte(data)); err != nil {
		return "", vars.ErrInternalServerError
	}

	token := hex.EncodeToString(h.Sum(nil)) + ":" + strconv.FormatInt(expires.Unix(), 10)

	return token, nil
}

func (tk *CSRFTokenManager) Verify(user_id uint64, cookie string, token string) (bool, error) {
	tokenData := strings.Split(token, ":")

	if len(tokenData) != 2 {
		return false, vars.ErrBadCSRF
	}

	tokenExp, err := strconv.ParseInt(tokenData[1], 10, 64)

	if err != nil {
		return false, vars.ErrInternalServerError
	}

	if tokenExp < time.Now().Unix() {
		return false, vars.ErrBadCSRF
	}

	h := hmac.New(sha1.New, tk.Secret)
	data := fmt.Sprintf("%d:%s:%d", user_id, cookie, tokenExp)
	h.Write([]byte(data))

	expected := h.Sum(nil)
	got, err := hex.DecodeString(tokenData[0])

	if err != nil {
		return false, vars.ErrInternalServerError
	}

	return hmac.Equal(expected, got), nil
}
