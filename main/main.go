package main

import (
    "fmt"
    "net/http"
    "github.com/adducci/jumpcloud/password"
    "time"
    "os"
    "strings"
    "bufio"
    "strconv"
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

//user chooses port number from 1024 to 49151
func getPort() string {
	reader := bufio.NewReader(os.Stdin)
	var input string

	//continue prompting until valid port entered or q entered to quit
    for havePort := false; havePort == false ; {
    	//get input
    	fmt.Print("Enter port: ")
        input, _ = reader.ReadString('\n')
        input = strings.TrimSpace(input)

        //q for quit
        if input == "q" {
        	return ""
        }

        //check if valid
        port, err := strconv.ParseInt(input, 10, 16)
        if err != nil  || port < 1024 || port > 49151 {
    	    fmt.Println("Please enter valid port number from 1024 to 49151 or enter q to quit")
        } else {
        	havePort = true
        }
    }

    return input
}

/*
Launch HTTP Server that listens to requests for /hash on provided port
*/
func main () {
	//get port
    port := getPort()
    if port == "" {
    	return
    }
    
    //listen on port
	var h hashHandler
	http.Handle("/hash", h)
    http.ListenAndServe(":" + port, nil)	
}
