package handler

import (
	"fmt"
	"net/http"
)

type Auth struct{}

func (auth *Auth) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Logging in")
}
