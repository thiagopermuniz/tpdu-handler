package main

import (
	"fmt"
	"net"
	"os"
)

const (
	SERVER_ADDRESS = "localhost:8583"
)

func main() {
	// Establish a TCP connection to the server
	conn, err := net.Dial("tcp", SERVER_ADDRESS)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	// Sample TPDU header + message
	tpdu := []byte{0x60, 0x00, 0x01, 0x00, 0x02}
	message := "Hello, Server with TPDU!"
	byteArray := []byte(message)
	fullMessage := append(tpdu, byteArray...)

	// Send the full message (TPDU + actual message) to the server
	_, err = conn.Write(fullMessage)
	if err != nil {
		fmt.Println("Error sending message:", err.Error())
		return
	}

	fmt.Println("Message with TPDU sent!")
}
