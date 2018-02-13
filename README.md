# jumpcloud

COMPILATION INSTRUCTIONS:
    1. Move to 'main' directory
    2. Run 'go build'
    3. Run 'main'
    4. Enter port number (int from 1024 to 49151) to run or 'q' to quit 

FILES:
    1. main 
        a. main.go
    2. password
        a. password.go
    3. httpserver
        a. myServer.go
        b. hashHandler.go
        c. statistics.go
        d. idMap.go

DESCRIPTION OF CODE AND LAYOUT:
    main
        main.go - prompts user for port number, launches server
    password
        password.go - short library for encrypting passwords using sha512 hash and base64 encoding
    httpserver
        myServer.go - implements version of http.Server unique to fit specifications of this program 
        hashHandler.go - defines server behavior, supports POST /hash, GET /hash{:id}, GET /stats, and /shutdown 
        statistics.go - implements a statistics object and related functions
        idMap.go - implements object that generates a unique integer id for an object and stores it in a map

DETAILED DESCRIPTION OF hashHandler:
    hashHandler.go
        Global Variables
            var hashes = IdMap{m : make(map[string]string)}
                stores the identifier : hash pairs processed by the server
            var st = Stats{0, -1}
                stores the total number of processed requests and average process request time
        Functions
	        func (h hashHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)
	            implements http.Handler interface
	            checks if r.Path and r.Method matches one of the supported requests
	            delegates writing response to matching function if it does
	            returns an appropriate error status code if not
	        func handleShutdown(w http.ResponseWriter, r *http.Request)
	            graceful shutdown
	        func handleStats(w http.ResponseWriter, r *http.Request)
	            responds with the servers current stats as json 
	        func handleGetHash(w http.ResponseWriter, r *http.Request)
	            responds with the hash identified in the request path
	        func handlePostHash(w http.ResponseWriter, r *http.Request) error
	            finds the password query from the request's form
	            immediately responds with an identifer that can be used to find that
	               password's hash in 5 seconds
	            returns error if no password query



