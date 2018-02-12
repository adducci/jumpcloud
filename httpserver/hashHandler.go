package httpserver

import (
    "fmt"
    "net/http"
    "time"
    "sync"
    "strconv"
    "github.com/adducci/jumpcloud/password"
)


//Type hashHandler implements the handler interface
type hashHandler struct {}

//Stores the hashes by their identifiers, enables locking
var hashes = struct {
	sync.RWMutex;
	m map[string]string;

} {
	m : make(map[string]string),
}

//Next id to return, incremental
var nextID int = 0





/*
Finds base64 encoded string of password with SHA512 hash
Lag for 5 seconds before storing it in the slice
*/
func computeHash(pw, id string) {
    time.Sleep(time.Second * 5)

    hash := password.Encrypt(pw)

    hashes.Lock()
    hashes.m[id] = hash
    hashes.Unlock()
}


/*
Returns the next id to use
*/
func getCurrentId(i chan int) {
   id := nextID
   nextID++
   i <- id
}

/*
Handle requests to /hash

Accept posts with a form field named password
Returns an identifer that can be used to retrieve the hash later

*/
func handleHash(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		//read in password value
		r.ParseForm()
		pw, ok := r.PostForm["password"]
		if ok {
			//get the next id to use
			i := make(chan int)
            go getCurrentId(i)
            id := <- i
            strID := strconv.Itoa(id)

			go computeHash(pw[0], strID)
			fmt.Fprint(w, strID)
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