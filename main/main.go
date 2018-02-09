package main

import (
    "fmt"
    "net/http"
    "github.com/adducci/jumpcloud/password"
)

type hashHandler struct {}


func (h hashHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//accept Post requests
	if r.Method == http.MethodPost {
		r.ParseForm()
		pw := r.PostForm["password"]
        fmt.Fprint(w, password.Encrypt(pw[0]))
	}
}

/*
Launch HTTP Server that listens to requests for /hash on port 8080
*/
func main () {
	var h hashHandler
	http.Handle("/hash", h)
    http.ListenAndServe(":8080", nil)	
}