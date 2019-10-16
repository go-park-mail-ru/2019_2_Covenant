package main

import (
	"2019_2_Covenant/pkg/user/delivery"
	_userRepo "2019_2_Covenant/pkg/user/repository"
	_sessionRepo "2019_2_Covenant/pkg/session/repository"
	"2019_2_Covenant/pkg/user/usecase"
	"github.com/labstack/echo"
	"log"
)

func main() {
	e := echo.New()

	userStorage := _userRepo.NewUserStorage()
	userUsecase := usecase.NewUserUsecase(userStorage)
	sessionStorage := _sessionRepo.NewSessionStorage()

	delivery.NewUserHandler(e, userUsecase, sessionStorage)

	log.Fatal(e.Start(":8000"))
}
