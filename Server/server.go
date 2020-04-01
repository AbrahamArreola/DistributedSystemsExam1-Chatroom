package main

import (
	"fmt"
	"net"
)

var clientList = make([]net.Conn, 0)
var data = make([]byte, 1024)

//Initialize listener port
func initServer(remote string) {
	fmt.Println("Initiating server... (<ENTER> to stop)")

	listener, err := net.Listen("tcp", remote)
	if err != nil {
		fmt.Printf("Error when listen: %s, Err: %s\n", remote, err)
		listener.Close()
		return
	}

	//Loop waiting for user data
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting client: ", err.Error())
			connection.Close()
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
		responseMessage = string(data[:length])
		userResponse := userName + " sent: " + responseMessage
		fmt.Println(userResponse)
		notifyUsers(connection, userResponse)
	}
}

func disconnectUser(connection net.Conn, userName string) {
	for index, con := range clientList {
		if con.RemoteAddr() == connection.RemoteAddr() {
			disconnectMessage := userName + " has left the room."
			fmt.Println(disconnectMessage)
			clientList = append(clientList[:index], clientList[index+1:]...)
			notifyUsers(connection, disconnectMessage)
		}
	}
}

func notifyUsers(connection net.Conn, message string) {
	for _, con := range clientList {
		if con.RemoteAddr() != connection.RemoteAddr() {
			con.Write([]byte(message))
		}
	}
}

func main() {
	var (
		stopInput string
		host      = "127.0.0.1"
		port      = "9000"
		remote    = host + ":" + port
	)

	go initServer(remote)
	fmt.Scanln(&stopInput)
}
