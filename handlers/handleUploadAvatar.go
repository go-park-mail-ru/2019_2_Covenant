package handlers

import "net/http"

func HandleUploadAvatar(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Upload avatar"))
}