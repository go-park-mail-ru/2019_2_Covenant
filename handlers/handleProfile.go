package server

import "net/http"

func HandleProfile(w http.ResponseWriter, r *http.Request) {
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