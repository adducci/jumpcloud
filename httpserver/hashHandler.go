package httpserver

import (
    "fmt"
    "net/http"
    "time"
    "sync"
    "strconv"
    "regexp"
    "encoding/json"
    "errors"
    "github.com/adducci/jumpcloud/password"
)


/*******************
TYPES
*******************/

//Type hashHandler implements the http.handler interface
type hashHandler struct {}


/*******************
GLOBAL VARIABLES
*******************/
//Stores the hashes by their identifiers, enables locking
var hashes = struct {
	sync.RWMutex;
	m map[string]string;

} {
	m : make(map[string]string),
}

//Stores the statitics on total requests and average time in ms to process them
var stats = struct {
	Total float64 "json:\"total\"";
	Average float64 "json:\"average\"";
} {
	Total : 0,
	Average : -1,
}

//Next id to return, incremental, equal to total number of requests
var nextID int = 0



/*******************
HELPER FUNCTIONS
*******************/

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
Client error: method not allowed
*/
func respondMethodNotAllowed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
    fmt.Fprint(w, "405 method not allowed")
}


/*
Signals that a password hash request has completed in time pt

Increases the total processes by 1, and adjusts the average time
according to the process time
*/
func updateStatistics(pt time.Duration) {
    ms := pt.Seconds() * 1e3

    if stats.Average < 0 {
    	//negative average implies that no requests have been processed yet
    	stats.Average = ms
    } else {
    	new_avg := (stats.Total * stats.Average + ms )/ (stats.Total + 1);
    	stats.Average = new_avg	
    }

    stats.Total++ 
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
        go getCurrentId(i)
        id := <- i
        strID := strconv.Itoa(id)
        
        //compute hash in background
		go computeHash(pw[0], strID)

		//immediately return identifier
		fmt.Fprint(w, strID)
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

Returns the hashes password for identifier id
*/
func handleGetHash(w http.ResponseWriter, r *http.Request) {
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
Handle requests to /stats
*/
func handleStats(w http.ResponseWriter, r *http.Request) {
	//encode the stats to json
    j, err := json.Marshal(stats)
  	if err != nil {
  		log.Println(err)
  	}

  	//return the encoded stats object
  	fmt.Fprint(w, string(j))
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
              updateStatistics(time.Since(start))
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