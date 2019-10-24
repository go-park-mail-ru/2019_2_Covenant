package main

import (
	_sessionRepo "2019_2_Covenant/internal/session/repository"
	_sessionUsecase "2019_2_Covenant/internal/session/usecase"
	"2019_2_Covenant/internal/user/delivery"
	_userRepo "2019_2_Covenant/internal/user/repository"
	_userUsecase "2019_2_Covenant/internal/user/usecase"
	"github.com/labstack/echo/v4"
	"log"

	_ "2019_2_Covenant/docs"
	"github.com/swaggo/echo-swagger"
)

// @title Covenant API
// @version 1.0
// @description Covenant backend server
// @BasePath /api/v1
func main() {
	e := echo.New()
	e.GET("/docs/*", echoSwagger.WrapHandler)

	userStorage := _userRepo.NewUserStorage()
	userUsecase := _userUsecase.NewUserUsecase(userStorage)
	sessionStorage := _sessionRepo.NewSessionStorage()
	sessionUsecase := _sessionUsecase.NewSessionUsecase(sessionStorage)

	handler := delivery.NewUserHandler(userUsecase, sessionUsecase)
	handler.Configure(e)

	log.Fatal(e.Start(":8000"))
}
