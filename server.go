package main

import (
	. "./handlers"
	. "./storage"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()

	api := &UsersHandler{
		Store: NewUserStore(),
		Session: NewSessionStore(),
	}

	e.POST("/api/v1/signup", api.SignUp)
	e.POST("/api/v1/signin", api.SignIn)

	log.Fatal(e.Start(":8000"))
}
