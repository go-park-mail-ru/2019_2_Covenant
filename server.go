package main

import (
	"2019_2_Covenant/pkg/user/delivery"
	"2019_2_Covenant/pkg/user/repository"
	"2019_2_Covenant/pkg/user/usecase"
	"github.com/labstack/echo"
	"log"
)

func main() {
	e := echo.New()

	userStorage := repository.NewUserStorage()
	userUsecase := usecase.NewUserUsecase(userStorage)
	delivery.NewUserHandler(e, userUsecase)

	log.Fatal(e.Start(":8000"))

	//
	//api := &UsersHandler{
	//	Store: NewUserStore(),
	//	Session: NewSessionStore(),
	//}
	//
	//e.POST("/api/v1/signup", api.SignUp)
	//e.POST("/api/v1/signin", api.SignIn)
	//
	//log.Fatal(e.Start(":8000"))
}
