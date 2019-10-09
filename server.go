package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	. "./handlers"
	. "./storage"
)


func main() {
	r := mux.NewRouter()

	api := &UsersHandler{
		Store: NewUserStore(),
		Session: NewSessionStore(),
	}

	//r.HandleFunc("/", ____)
	r.HandleFunc("/login", api.SignIn).Methods("POST")
	r.HandleFunc("/logout", api.SignOut).Methods("GET")
	r.HandleFunc("/signup", api.SignUp).Methods("POST")
	r.HandleFunc("/profile", api.GetProfile).Methods("GET")
	r.HandleFunc("/profile", api.PostProfile).Methods("POST")
	//r.HandleFunc("/upload/avatar", ____)

	log.Println("start serving :8080")
	http.ListenAndServe(":8080", r)
}
