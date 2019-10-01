package server

import "net/http"

func HandleSignUp(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Sign Up"))
}