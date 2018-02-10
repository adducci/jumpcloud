package httpserver

import (
    "net/http"
    // "fmt"
    "context"
    "sync"
)

var s http.Server

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
Shutsdown the server with no cancel or deadline
*/
func shutdownMyServer() {
	s.Shutdown(context.Background())
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
        s.ListenAndServe()
    }()

    wait.Wait()
}