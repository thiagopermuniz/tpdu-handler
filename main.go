package main

import (
	"bytes"
	"encoding/hex"
	_ "encoding/hex"
	"fmt"
	"net"
	"time"
)

const (
	PORT            = ":8583"
	PROXY_SERVER_01 = "localhost:8801"
	PROXY_SERVER_02 = "localhost:8802"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}
	buffer = buffer[:n]
	for _, b := range buffer {
		fmt.Printf("%02x ", b)
	}
	tpdu := buffer[:5]
	//origin := tpdu[1:3]
	dest := tpdu[3:5]
	byteData, err := hex.DecodeString(string(buffer))
	fmt.Println(byteData)
	switch {
	case bytes.Equal(dest, []byte{0x00, 0x01}):
		fmt.Println("01 dest message")
		response, err := proxyMessage(buffer, PROXY_SERVER_01)
		if err != nil {
			fmt.Println("Error proxying message:", err.Error())
			return
		}
		_, err = conn.Write(response)
		if err != nil {
			fmt.Println("Error sending response to client:", err.Error())
		}
	case bytes.Equal(dest, []byte{0x00, 0x02}):
		fmt.Println("02 dest message")
		response, err := proxyMessage(buffer, PROXY_SERVER_02)
		if err != nil {
			fmt.Println("Error proxying message:", err.Error())
			return
		}

		_, err = conn.Write(response)
		if err != nil {
			fmt.Println("Error sending response to client:", err.Error())
		}
	default:
		fmt.Println("Unknown origin.")
	}

}

func proxyMessage(message []byte, targetHost string) ([]byte, error) {
	// Connect to the target proxy server
	proxyConn, err := net.Dial("tcp", targetHost)
	if err != nil {
		return nil, err
	}
	defer proxyConn.Close()

	// Send the entire message to the proxy server
	_, err = proxyConn.Write(message)
	if err != nil {
		return nil, err
	}

	// Now wait for the response from the proxy server
	response := make([]byte, 1024)
	n, err := proxyConn.Read(response)
	if err != nil {
		return nil, err
	}
	return response[:n], nil
}

func main() {
	for {
		listener, err := net.Listen("tcp", PORT)
		if err != nil {
			fmt.Println("Error listening:", err.Error())
			time.Sleep(5 * time.Second) // wait for 5 seconds before trying again
			continue                    // jump back to the start of the loop
		}

		fmt.Println("Listening on " + PORT)
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error accepting:", err.Error())
				break // this will exit the inner loop, close the listener, and try to re-listen
			}
			go handleConnection(conn)
		}

		listener.Close() // Close the listener if we break out of the inner loop
	}
}
