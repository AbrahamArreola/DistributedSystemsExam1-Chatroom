package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

//DataFile is...
type DataFile struct {
	Name     string
	ByteFile []byte
}

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
			sendDataType(connection, true)
			messageToServer(connection)

		case 2:
			sendDataType(connection, false)
			fileToServer(connection)

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

//Function to notify the server whether the message is a string or a file
func sendDataType(connection net.Conn, isMessage bool) {
	if isMessage {
		writeData = []byte("1")
	} else {
		writeData = []byte("0")
	}
	input, err := connection.Write([]byte(writeData))
	if err != nil {
		fmt.Printf("Error when send to server: %d\n", input)
		return
	}
}

//Function to write data to the server to send messages
func messageToServer(connection net.Conn) {
	fmt.Print("Write your message: ")
	writeData, _, _ = scan.ReadLine()
	input, err := connection.Write([]byte(writeData))
	if err != nil {
		fmt.Printf("Error when send to server: %d\n", input)
		return
	}
	allMessages += "You sent: " + string(writeData) + "\n"
}

//Function to send a file to the server
func fileToServer(connection net.Conn) {
	var dataFile DataFile

	fmt.Print("Write the route of the file: ")
	writeData, _, _ = scan.ReadLine()
	route := string(writeData)
	fmt.Println(route)

	//Get file info
	fileInfo, err := os.Stat(route)
	//Open the file
	file, err := os.Open(route)
	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	dataFile.Name = fileInfo.Name()
	dataFile.ByteFile = data
	file.Close()

	//Send file to server
	err = gob.NewEncoder(connection).Encode(dataFile)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	var (
		host   = "127.0.0.1"
		port   = "9000"
		remote = host + ":" + port
	)

	joinUser(remote)
}
