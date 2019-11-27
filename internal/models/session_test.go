package models

import (
	"2019_2_Covenant/tools/vars"
	"fmt"
	"testing"
	"time"
)
var token string

func TestCSRFTokenManager_Create(t *testing.T) {
	CSRFTokenManager := NewCSRFTokenManager("Covenant")

	var err error

	token, err = CSRFTokenManager.Create(uint64(1), "covenant", time.Now().Add(24*time.Hour))
	if token == "" || err != nil {
		fmt.Println("Token: expected not nil, got ", token)
		fmt.Println("Error: expected nil, got ", err)
		t.Fail()
	}
}

func TestCSRFTokenManager_Verify(t *testing.T) {
	CSRFTokenManager := NewCSRFTokenManager("Covenant")

	t.Run("Test OK", func(t1 *testing.T){
		ok, err := CSRFTokenManager.Verify(uint64(1), "covenant", token)

		if !ok && err != nil {
			fmt.Println("Expected true, got ", ok)
			fmt.Println("Error: expected nil, got ", err)
			t1.Fail()
		}
	})

	t.Run("Error of token", func(t2 *testing.T){
		wrongToken := "wrong token"
		ok, err := CSRFTokenManager.Verify(uint64(1), "covenant", wrongToken)

		if ok && err != vars.ErrBadCSRF {
			fmt.Println("Expected false, got ", ok)
			fmt.Println("Error: expected csrf error, got ", err)
			t2.Fail()
		}
	})

	t.Run("Error of expiring", func(t3 *testing.T){
		newToken, _ := CSRFTokenManager.Create(uint64(1), "covenant", time.Date(2018, 1,1,0,0,0,0, time.Local))
		ok, err := CSRFTokenManager.Verify(uint64(1), "covenant", newToken)

		if ok && err != vars.ErrBadCSRF {
			fmt.Println("Expected false, got ", ok)
			fmt.Println("Error: expected csrf error, got ", err)
			t3.Fail()
		}
	})

}
