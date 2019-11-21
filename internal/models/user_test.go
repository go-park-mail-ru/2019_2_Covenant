package models

import (
	"fmt"
	"testing"
)

func TestUser_BeforeStore(t *testing.T) {
	user := &User{
		ID: 1, Nickname: "nick", Email: "e@mail.ru", PlainPassword: "qwerty", Avatar: "path", Role: 0, Access: 0,
	}

	err := user.BeforeStore()
	if err != nil {
		fmt.Println("Error: expected nil, got", err)
		t.Fail()
	}
}

func TestUser_Verify(t *testing.T) {
	user := &User{
		ID: 1, Nickname: "nick", Email: "e@mail.ru", PlainPassword: "qwerty", Avatar: "path", Role: 0, Access: 0,
	}

	_ = user.BeforeStore()

	ok := user.Verify("qwerty")
	if ok != true {
		t.Fail()
	}
}
