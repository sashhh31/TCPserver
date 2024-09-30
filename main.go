package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

const (
	port         = ":2504"
	readTimeout  = 10 * time.Second // Timeout duration for reading data
	writeTimeout = 10 * time.Second // Timeout duration for writing data
)

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error while listening:", err)
		return
	}
	defer listener.Close()
	fmt.Printf("Server is listening on port %s...\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error while accepting connection:", err)
			continue // Continue accepting new connections instead of returning
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		// Set a read timeout
		conn.SetReadDeadline(time.Now().Add(readTimeout))

		// Read message from client
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}

		trimmedMessage := strings.TrimSpace(message)
		fmt.Printf("Received: %s\n", trimmedMessage)

		if strings.HasPrefix(trimmedMessage, "GET") {
			response := "HTTP/1.1 200 OK\r\n" +
				"Content-Type: text/plain\r\n" +
				"Connection: close\r\n\r\n" +
				"Hello from TCP server!\n"

			// Set a write timeout
			conn.SetWriteDeadline(time.Now().Add(writeTimeout))

			_, err = conn.Write([]byte(response))
			if err != nil {
				fmt.Println("Error writing to connection:", err)
				return
			}
			return
		}

		// Regular message handling
		conn.SetWriteDeadline(time.Now().Add(writeTimeout))
		_, err = conn.Write([]byte("Message received.\n"))
		if err != nil {
			fmt.Println("Error writing to connection:", err)
			return
		}
	}
}
