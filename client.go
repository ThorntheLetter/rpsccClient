package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
)

func main() {
	// user's program io setup
	userName := os.Args[1]
	userProgram := os.Args[2]
	userArgs := os.Args[3:]
	userCmd := exec.Command(userProgram, userArgs...)
	userOutPipe, err := userCmd.StdoutPipe()
	if err != nil {
		fmt.Println(err.Error())
	}
	userOut := bufio.NewReader(userOutPipe)

	userInPipe, err := userCmd.StdinPipe()
	if err != nil {
		fmt.Println(err.Error())
	}

	err = userCmd.Start()
	if err != nil {
		fmt.Println(err)
	}

	// network setup
	connection, err := net.Dial("tcp", "localhost:7633")
	if err != nil {
		fmt.Println(err.Error()) // this error is connection problems
	}
	defer connection.Close()
	netReader := bufio.NewReader(connection)
	connection.Write([]byte(userName + "\n")) // send name

	for {
		recieved, err := netReader.ReadString('\n')
		if err != nil {
			fmt.Println(err.Error())
		}
		if recieved == "6\n" {
			fmt.Printf("%s", recieved)

			_, err = userInPipe.Write([]byte(recieved))
			if err != nil {
				fmt.Println(err.Error())
			}

			for i := 0; i < 100; i += 1 {
				choice, err := userOut.ReadString('\n')
				if err != nil {
					fmt.Println(err.Error())
				}
				fmt.Printf(choice)
				_, err = connection.Write([]byte(choice))
				if err != nil {
					fmt.Println(err.Error())
				}

				opChoice, err := netReader.ReadString('\n')
				if err != nil {
					fmt.Println(err.Error())
				}
				fmt.Printf(opChoice)
				_, err = userInPipe.Write([]byte(opChoice))
				if err != nil {
					fmt.Println(err)
				}

			}
		}
	}
}
