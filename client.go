package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
)

func userErrorHandler(e error) { // this gives the user the error because i have no idea what would actually cause these
	fmt.Println(e.Error())
	fmt.Println("something probably went wrong with your program, check above for possible details")
	os.Exit(1)
}

func disconnected() {
	fmt.Println("you probably have been disconnected")
	os.Exit(2) // still on the fence about whether this should exit or break
}

func main() {
	// user's program io setup
	userName := os.Args[1]

	if userName == "cody" {
		connection, err := net.Dial("tcp", "rpscc.thebelligerentone.com:3489")
		if err != nil {
			fmt.Println("unable to connect to server") // this error is connection problems
			os.Exit(1)
		}
		r := bufio.NewReader(connection)
		for {
			q, err := r.ReadString('\n')
			if err != nil {
				os.Exit(0)
			}
			fmt.Printf(q)
		}

	}

	userProgram := os.Args[2]
	userArgs := os.Args[3:]
	userCmd := exec.Command(userProgram, userArgs...)
	userOutPipe, err := userCmd.StdoutPipe()
	if err != nil {
		userErrorHandler(err)
	}
	userOut := bufio.NewReader(userOutPipe)

	userInPipe, err := userCmd.StdinPipe()
	if err != nil {
		userErrorHandler(err)
	}

	err = userCmd.Start()
	if err != nil {
		userErrorHandler(err)
	}

	// network setup
	connection, err := net.Dial("tcp", "rpscc.thebelligerentone.com:7633")
	if err != nil {
		fmt.Println("unable to connect to server") // this error is connection problems
		os.Exit(1)
	}
	defer connection.Close()
	netReader := bufio.NewReader(connection)
	_, err = connection.Write([]byte(userName + "\n")) // send name
	if err != nil {
		disconnected()
	}

	// the loop
	for {
		recieved, err := netReader.ReadString('\n')
		if err != nil {
			disconnected()
		}
		if recieved == "6\n" {
			fmt.Printf("%s", recieved)

			_, err = userInPipe.Write([]byte(recieved))
			if err != nil {
				userErrorHandler(err)
			}

			for i := 0; i < 100; i += 1 {
				choice, err := userOut.ReadString('\n') // read user choice
				if err != nil {
					userErrorHandler(err)
				}
				fmt.Printf(choice)                        // print user choice
				_, err = connection.Write([]byte(choice)) // send user choice
				if err != nil {
					disconnected()
				}

				opChoice, err := netReader.ReadString('\n') // recieve opponents choice
				if err != nil {
					disconnected()
				}

				fmt.Printf(opChoice)
				_, err = userInPipe.Write([]byte(opChoice)) // give user the opponents choice
				if err != nil {
					userErrorHandler(err)
				}

				forfeit := false
				switch choice {
				case "1\n", "2\n", "3\n", "4\n", "5\n":
				default:
					forfeit = true
				}
				if opChoice == "0\n" || forfeit { // end loop if user forfeited
					break
				}
			}
		}
	}
}
