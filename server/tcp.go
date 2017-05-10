package main

import (
	"github.com/alinvasile/avcfg/common"
	"io"
	"log"
	"fmt"
	"net"
	"os"
	"github.com/patrickmn/go-cache"
)

var (
	connections []net.Conn
)

func startTcpServer(propertyCache *cache.Cache, port string){
	l, err := net.Listen("tcp", "0.0.0.0"+":"+port)
	if err != nil {
		fmt.Println("Error listening on tcp:", err.Error())
		os.Exit(1)
	}

	defer l.Close()
	fmt.Println("TCP Listening on " + port)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Save connection
		connections = append(connections, conn)
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	for {
		msg, err := common.ReadMessage(conn)
		if err != nil {
			if err == io.EOF {
				// Close the connection when you're done with it.
				removeConn(conn)
				conn.Close()
				return
			}
			log.Println(err)
			return
		}
		fmt.Printf("Message Received: %s\n", msg)
		broadcast(conn, msg)
	}
}

func removeConn(conn net.Conn) {
	var i int
	for i = range connections {
		if connections[i] == conn {
			break
		}
	}
	connections = append(connections[:i], connections[i+1:]...)
}

func broadcast(conn net.Conn, msg string) {
	for i := range connections {
		if connections[i] != conn {
			err := common.WriteMessage(connections[i], msg)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
