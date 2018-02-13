package httpserver

import (
    "net/http"
    "context"
    "sync"
    "log"
)


/**************
GLOBAL VARIABLES
***************/
//The server
var s http.Server



/****************
FUNCTIONS
****************/
/*
initalizes an http server handled by hashHandler
*/
func makeServer(port string) {
	h := hashHandler{}
	s = http.Server{
		Addr: "localhost:" + port,
		Handler: h,
	}
}

/*
Shuts down the server with no cancel or deadline
*/
func shutdownMyServer() {
	err := s.Shutdown(context.Background())
    if err != nil {
        log.Fatal(err)
    }
}

/*
Initalizes and runs a server that listens on the given port and 
serves the hashHandler on /hash posts
*/
func Run(port string) {
    //initalize server
    makeServer(port)

    //run the server, and wait for it to shutdown
    var wait sync.WaitGroup
	wait.Add(1)

	go func () { 
        defer wait.Done()
        err := s.ListenAndServe()
        if err != nil {
            log.Fatal(err)
        }
    }()


    wait.Wait()
}