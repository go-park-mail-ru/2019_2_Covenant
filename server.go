package main

import (
	handlers "./handlers"
	"net/http"
)

type User struct {
	ID       uint64 `json:"id,uint64"`
	Username string `json:"username,string"`
	Email    string `json:"email,string"`
	Password string `json:"-"`
	Avatar   string `json:"avatar,string"`
}

type UserInput struct {
	Email    string `json:"email,string"`
	Password string `json:"password,string"`
}

func main() {
	http.HandleFunc("/", handlers.HandleMain)
	http.HandleFunc("/login", handlers.HandleLogin)
	http.HandleFunc("/logout", handlers.HandleLogout)
	http.HandleFunc("/signup", handlers.HandleSignUp)
	http.HandleFunc("/profile", handlers.HandleProfile)
	http.HandleFunc("/upload/avatar", handlers.HandleUploadAvatar)

	http.ListenAndServe(":3000", nil)
}
