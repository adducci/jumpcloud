package main

import (
	"bufio"
	"fmt"
	"github.com/adducci/jumpcloud/httpserver"
	"os"
	"strconv"
	"strings"
)

//user chooses port number from 1024 to 49151
func getPort() string {
	reader := bufio.NewReader(os.Stdin)
	var input string

	//continue prompting until valid port entered or q entered to quit
	for havePort := false; havePort == false; {
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
		if err != nil || port < 1024 || port > 49151 {
			fmt.Println("Please enter valid port number from 1024 to 49151 or enter q to quit")
		} else {
			havePort = true
		}
	}

	return input
}

/*
Launch HTTP Server that listens to requests on provided port
*/
func main() {
	//get port
	port := getPort()
	if port == "" {
		return
	}

	//run myServer on given Port
	httpserver.Run(port)
}
