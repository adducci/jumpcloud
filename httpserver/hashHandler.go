package httpserver

import (
    "fmt"
    "net/http"
    "time"
    "sync"
    "strconv"
    "regexp"
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
Lag for 5 seconds before storing it in the map
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
Handle post requests to /hash

Accept posts with a form field named password
Returns an identifer that can be used to retrieve the hash later

*/
func postHash(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
    	w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "400 needs form")
    }

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
		fmt.Fprint(w, "400 needs password query")
	} 
}


/*
Handle get requests to /hash/{:id}

Returns the hashes password for identifier id
*/
func getHash(w http.ResponseWriter, r *http.Request) {
	//get id from path 
	reg, _ := regexp.Compile("[0-9]+")
	id := reg.Find([]byte(r.URL.Path))

    //get hash from hashes
	hashes.RLock()
	hash := hashes.m[string(id)]
	hashes.RUnlock()

    //return hash 
	if hash != "" {
        fmt.Fprint(w, hash)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "404 hash not found")
	}
}

/*
Handle requests to /shutdown
*/
func handleShutdown(w http.ResponseWriter, r *http.Request) {
    shutdownMyServer()
}

/*
Client error: method not allowed
*/
func respondMethodNotAllowed(w http.ResponseWriter) {
	 w.WriteHeader(http.StatusMethodNotAllowed)
     fmt.Fprint(w, "405 method not allowed")
}


/*
Handle all requests to the server
/hash, /hash{:id} and /shutdown only valid paths
*/
func (h hashHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path 

	if match, _ := regexp.MatchString("^/?hash$", path); match {
		// requests to /hash should be posts
		if r.Method == http.MethodPost {
            postHash(w, r)
		} else {
			//otherwise method not allowed
            respondMethodNotAllowed(w)
		}
	} else if match, _ := regexp.MatchString("^/?hash/[0-9]+", path); match {
		//requests to /hash/{:id} should be get
		if r.Method == http.MethodGet {
			getHash(w, r)
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