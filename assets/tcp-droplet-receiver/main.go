package main

import (
	"bytes"
	"flag"
	"fmt"
	"net"
	"os"
)

const (
	DEFAULT_ADDRESS   = "0.0.0.0:3333"
	CONN_TYPE         = "tcp"
	DEFAULT_SERVER_ID = "droplet_server"
)

var serverAddress = flag.String(
	"address",
	DEFAULT_ADDRESS,
	"The host:port that the server is bound to.",
)

var serverId = flag.String(
	"serverId",
	DEFAULT_SERVER_ID,
	"The Server id that is echoed back for each message.",
)

func main() {
	flag.Parse()
	// Listen for incoming connections.
	listener, err := net.Listen(CONN_TYPE, *serverAddress)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer listener.Close()
	fmt.Printf("%s:Listening on %s\n", *serverId, *serverAddress)
	for {
		// Listen for an incoming connection.
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Close the connection when you're done with it.
	defer conn.Close()
	// Make a buffer to hold incoming data.
	buff := make([]byte, 1024)
	// Continue to receive the data forever...
	for {
		// Read the incoming connection into the buffer.
		readBytes, err := conn.Read(buff)
		if err != nil {
			fmt.Println("Closing connection:", err.Error())
			return
		}
		var writeBuffer bytes.Buffer
		writeBuffer.WriteString(*serverId)
		writeBuffer.WriteString(":")
		writeBuffer.Write(buff[0:readBytes])
		fmt.Print(writeBuffer.String())
		_, err = conn.Write(writeBuffer.Bytes())
		if err != nil {
			fmt.Println("Closing connection:", err.Error())
			return
		}
	}
}
