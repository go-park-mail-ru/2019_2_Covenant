package server

import "net/http"

func HandleMain(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Main page"))
}
