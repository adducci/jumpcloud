package main

import (
    "fmt"
    "net/http"
    "github.com/adducci/jumpcloud/password"
    "time"
)

type hashHandler struct {}


/*
Accept posts with a form field named password
Respond with base64 encoded string of password with SHA512 hash
Lag for 5 seconds before responding
*/
func (h hashHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//respon to only Post requests
	if r.Method == http.MethodPost {
		//encrypt given password if given
		r.ParseForm()
		pw, ok := r.Form["password"]
		if ok {
            fmt.Fprint(w, password.Encrypt(pw[0]))
		}
	}
	time.Sleep(time.Second * 5)
}

/*
Launch HTTP Server that listens to requests for /hash on port 8080
*/
func main () {
	var h hashHandler
	http.Handle("/hash", h)
    http.ListenAndServe(":8080", nil)	
}