package main

import (
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Main page"))
	})
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", handleLogout)
	http.HandleFunc("/signup", handleSignUp)
	http.HandleFunc("/profile", handleProfile)
	http.HandleFunc("/upload/avatar", handleUploadAvatar)

	http.ListenAndServe(":3000", nil)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Login"))
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Logout"))
}

func handleSignUp(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Sign Up"))
}

func handleProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		handleProfileGet(w, r)
	} else if r.Method == http.MethodPost {
		handleProfilePost(w, r)
	}
}

func handleProfileGet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Profile Get"))
}

func handleProfilePost(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Profile Post"))
}

func handleUploadAvatar(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Upload avatar"))
}
