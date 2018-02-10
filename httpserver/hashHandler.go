package httpserver

import (
    "fmt"
    "net/http"
    "time"
    "github.com/adducci/jumpcloud/password"
)


/*
Type hashHandler implements the handler interface
*/
type hashHandler struct {}


/*
Handle requests to /hash

Accept posts with a form field named password
Respond with base64 encoded string of password with SHA512 hash
Lag for 5 seconds before responding
*/
func handleHash(w http.ResponseWriter, r *http.Request) {
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
Handle requests to /shutdown
*/
func handleShutdown(w http.ResponseWriter, r *http.Request) {
    shutdownMyServer()
}


/*
Handle all requests to the server
Only accepts /hash and /shutdown paths
*/
func (h hashHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "hash" || path == "/hash" {
        handleHash(w, r)
	} else if path == "shutdown" || path == "/shutdown" {
        handleShutdown(w, r)
	} else {
        //no other paths accpeted 
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "404 not found")
	}
}