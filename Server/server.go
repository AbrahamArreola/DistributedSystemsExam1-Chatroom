package main

import (
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

var clientList = make([]net.Conn, 0)
var data = make([]byte, 1024)
var allMessages string

//Initialize listener port
func initServer(remote string) {
	fmt.Println("Running server...")

	listener, err := net.Listen("tcp", remote)
	if err != nil {
		fmt.Printf("Error when listen: %s, Err: %s\n", remote, err)
		return
	}

	go optionMenu()

	//Loop waiting for user data
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting client: ", err.Error())
			return
		}
		clientList = append(clientList, connection)
		go manageUsers(connection)
	}
}

//Managing user data
func manageUsers(connection net.Conn) {
	fmt.Println("New user connected in: ", connection.RemoteAddr())

	//First it's created a reader to get user data from the port
	length, err := connection.Read(data)
	if err != nil {
		fmt.Printf("Client %s quit.\n", connection.RemoteAddr())
		connection.Close()
		disconnectUser(connection, connection.RemoteAddr().String())
		return
	}

	//Then it gets the username
	userName := string(data[:length])
	welcomeMessage := userName + " entered the chat."
	allMessages += welcomeMessage + "\n"
	notifyUsers(connection, welcomeMessage)

	//Finally it receives and sends messages
	//Variable for message response
	var responseMessage string
	//Variable for file response
	for {
		length, err := connection.Read(data)
		if err != nil {
			connection.Close()
			disconnectUser(connection, userName)
			return
		}

		dataType := string(data[:length])
		if dataType == "1" {
			length, err = connection.Read(data)
			if err != nil {
				connection.Close()
				disconnectUser(connection, userName)
				return
			}

			responseMessage = string(data[:length])
			userResponse := userName + " sent: " + responseMessage
			allMessages += userResponse + "\n"
			notifyUsers(connection, userResponse)

		} else {
			var dataFile DataFile
			err := gob.NewDecoder(connection).Decode(&dataFile)
			if err != nil {
				fmt.Println(err)
				return
			}
			userResponse := userName + " sent: " + dataFile.Name
			allMessages += userResponse + "\n"
			notifyUsers(connection, userResponse)
			saveFileSent(dataFile)
		}
	}
}

func saveFileSent(dataFile DataFile) {
	dir := "files\\" + dataFile.Name

	err := ioutil.WriteFile(dir, []byte(dataFile.ByteFile), 0666)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("New file created!")
	}
}

//Search function that iterates over clientList to disconnect certain user
func disconnectUser(connection net.Conn, userName string) {
	for index, con := range clientList {
		if con.RemoteAddr() == connection.RemoteAddr() {
			disconnectMessage := userName + " has left the chat."
			allMessages += disconnectMessage + "\n"
			//fmt.Println(disconnectMessage)
			clientList = append(clientList[:index], clientList[index+1:]...)
			notifyUsers(connection, disconnectMessage)
		}
	}
}

//Function to send messages to all users connect except the one who sends them
func notifyUsers(connection net.Conn, message string) {
	for _, con := range clientList {
		if con.RemoteAddr() != connection.RemoteAddr() {
			con.Write([]byte(message))
		}
	}
}

//Option menu to select an option to do
func optionMenu() {
	var option int
	for {
		fmt.Println("Options:\n1.Show messages\n2. Backup messages\n3. Stop server")
		fmt.Scanln(&option)
		switch option {
		case 1:
			fmt.Printf("Option: %d\n***** All messages in chatroom*****\n%s", option, allMessages)
			fmt.Print("***********************************\n")

		case 2:
			err := ioutil.WriteFile("messages.txt", []byte(allMessages), 0666)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Messages backup done!")

		case 3:
			fmt.Println("Server stopped")
			os.Exit(0)
		}
	}
}

func main() {
	var (
		host   = "127.0.0.1"
		port   = "9000"
		remote = host + ":" + port
	)

	initServer(remote)
}
