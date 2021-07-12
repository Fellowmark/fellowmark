package routes

import (
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Logged in"))
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Signed up"))
}
