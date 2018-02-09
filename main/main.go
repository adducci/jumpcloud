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
	if r.Method == http.MethodPost {
		//read in password value
		r.ParseForm()
		pw, ok := r.PostForm["password"]
		if ok {
			//encrypt and return hashed password if given
            fmt.Fprint(w, password.Encrypt(pw[0]))
            time.Sleep(time.Second * 5)
		} else {
			//return error if not password given 
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "400 bad request - needs password query")
		}
	} else {
		//Post is only supported method
         w.WriteHeader(http.StatusMethodNotAllowed)
         fmt.Fprint(w, "<p> 405 method not allowed </p>")
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

func errorHandle(status int, message string) {

}