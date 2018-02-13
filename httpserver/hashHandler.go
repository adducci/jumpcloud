package httpserver

import (
	"errors"
	"fmt"
	"github.com/adducci/jumpcloud/password"
	"net/http"
	"regexp"
	"time"
)

/*******************
TYPES
*******************/

//Type hashHandler implements the http.handler interface
type hashHandler struct{}

/*******************
GLOBAL VARIABLES
*******************/
//Stores the hashes by their identifiers, enables locking
var hashes = IdMap{m: make(map[string]string)}

//Statistics on this server
var st = Stats{0, -1}

/*******************
HELPER FUNCTIONS
*******************/

/*
Finds base64 encoded string of password with SHA512 hash
Lag for 5 seconds before storing it in the map
*/
func computeHash(pw string, id int) {
	time.Sleep(time.Second * 5)

	hash := password.Encrypt(pw)

	hashes.WriteToMap(hash, id)
}

/*
Client error: method not allowed
*/
func respondMethodNotAllowed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprint(w, "405 method not allowed")
}

/*******************
HANDLE DIFFERENT PATH REQUESTS
*******************/

/*
Handle post requests to /hash

Accept posts with a form field named password
Returns an identifer that can be used to retrieve the hash later

*/
func handlePostHash(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		//return error if no form to parse
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "400 needs form")
		return errors.New("Post /hash needs form")
	}

	pw, ok := r.PostForm["password"]

	if ok {
		//get the next id to use
		i := make(chan int)
		go GetCurrentId(i)
		id := <-i

		//compute hash in background
		go computeHash(pw[0], id)

		//immediately return identifier
		fmt.Fprint(w, id)
	} else {
		//return error if not password given
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "400 needs password query")
		return errors.New("Post /hash needs password query")
	}

	return nil
}

/*
Handle get requests to /hash/{:id}

Returns the hashed password for identifier id
*/
func handleGetHash(w http.ResponseWriter, r *http.Request) {
	//get id from path
	reg, _ := regexp.Compile("[0-9]+")
	id := reg.Find([]byte(r.URL.Path))

	//get hash from hashes
	hash := hashes.ReadFromMap(string(id))

	//return hash
	if hash != "" {
		fmt.Fprint(w, hash)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "404 hash not found")
	}
}

/*
Handle requests to /stats
*/
func handleStats(w http.ResponseWriter, r *http.Request) {
	//encode the stats to json
	j := st.Encode()

	//return the encoded stats object
	fmt.Fprint(w, j)
}

/*
Handle requests to /shutdown
*/
func handleShutdown(w http.ResponseWriter, r *http.Request) {
	shutdownMyServer()
}

/*******************
PARSE THE INCOMING HTTP REQUESTS
*******************/

/*
Handle all requests to the server
/hash, /hash{:id}, /stats, and /shutdown only valid paths
Routes to the matched path if correct method
*/
func (h hashHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	path := r.URL.Path

	if match, _ := regexp.MatchString("^/?hash$", path); match {
		// requests to /hash should be posts
		if r.Method == http.MethodPost {
			err := handlePostHash(w, r)
			if err == nil {
				//update stats with successful posts' process time
				st.UpdateStatistics(time.Since(start))
			}
		} else {
			//otherwise method not allowed
			respondMethodNotAllowed(w)
		}
	} else if match, _ := regexp.MatchString("^/?hash/[0-9]+", path); match {
		//requests to /hash/{:id} should be get
		if r.Method == http.MethodGet {
			handleGetHash(w, r)
		} else {
			//otherwise method not allowed
			respondMethodNotAllowed(w)
		}
	} else if match, _ := regexp.MatchString("^/?stats$", path); match {
		// requests to /stats should be get
		if r.Method == http.MethodGet {
			handleStats(w, r)
		} else {
			//otherwise method not allowed
			respondMethodNotAllowed(w)
		}
	} else if match, _ := regexp.MatchString("^/?shutdown", path); match {
		//handles requests to shutdown
		handleShutdown(w, r)
	} else {
		//no other paths supported
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "400 invalid path")
	}
}
