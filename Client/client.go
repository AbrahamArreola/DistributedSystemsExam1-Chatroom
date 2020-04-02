package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

var writeData = make([]byte, 1024)
var readData = make([]byte, 1024)
var scan = bufio.NewReader(os.Stdin)
var allMessages string

//Function to open new client connection
func joinUser(remote string) {
	connection, err := net.Dial("tcp", remote)

	if err != nil {
		fmt.Println("Server not found.")
		return
	}

	fmt.Println("You've entered the chat.")
	allMessages += "You've entered the chat.\n"

	//First it's sent the username
	fmt.Print("Enter your name: ")
	fmt.Scanln(&writeData)
	input, err := connection.Write([]byte(writeData))
	if err != nil {
		fmt.Printf("Error when send to server: %d\n", input)
		return
	}

	//Then it's called concurrently readFromServer to listen for data from the server
	var option int
	go readFromServer(connection, option)

	//Loop to select an option to do
	for {
		fmt.Println("Options:\n1. send message\n2. send file\n3. show chatroom\n4. Exit chatroom")
		fmt.Scanln(&option)
		switch option {
		case 1:
			fmt.Printf("Option: %d\n", option)
			writeToServer(connection)

		case 3:
			fmt.Printf("Option: %d\n***** All messages since you entered*****\n%s", option, allMessages)
			fmt.Print("*****************************************\n")

		case 4:
			fmt.Println("You've left the chatroom")
			return

		default:
			fmt.Printf("There is no %d option\n", option)
		}
	}
}

//Function to read data from server while option different of 4 (exit program)
func readFromServer(connection net.Conn, option int) {
	for option != 4 {
		length, err := connection.Read(readData)
		if err != nil {
			fmt.Printf("Error when read from server. Error:%s\n", err)
			os.Exit(0)
		}
		fmt.Println(string(readData[:length]))
		allMessages += string(readData[:length]) + "\n"
	}
	connection.Close()
	return
}

//Function to write data to the server to send messages
func writeToServer(connection net.Conn) {
	fmt.Print("Write your message: ")
	writeData, _, _ = scan.ReadLine()
	input, err := connection.Write([]byte(writeData))
	if err != nil {
		fmt.Printf("Error when send to server: %d\n", input)
		return
	}
	allMessages += "You sent: " + string(writeData) + "\n"
}

func main() {
	var (
		host   = "127.0.0.1"
		port   = "9000"
		remote = host + ":" + port
	)

	joinUser(remote)
}
