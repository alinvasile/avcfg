package main

import (
	"flag"
	"net"
	"log"
	"github.com/alinvasile/avcfg/common"
	"io"
	"fmt"
	"os"
	"bufio"

)

// based on https://github.com/jonfk/golang-chat/blob/master/tcp/client-chat/main.go

var tcp_port_string = flag.String("tcp.port", "8087", "TCP Listen port")
var tcp_host_string = flag.String("tcp.host", "127.0.0.1", "TCP Listen host")

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	flag.Parse()

	host := *tcp_host_string
	port := *tcp_port_string

	tcpAddr, err := net.ResolveTCPAddr("tcp", host+":"+port)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to server through tcp.
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	go printOutput(conn)
	writeInput(conn)
}

func writeInput(conn *net.TCPConn) {

	// Read from stdin.
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter command: ")
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		err = common.WriteMessage(conn,text)
		if err != nil {
			log.Println(err)
		}
	}
}

func printOutput(conn *net.TCPConn) {
	for {
		msg, err := common.ReadMessage(conn)
		// Receiving EOF means that the connection has been closed
		if err == io.EOF {
			// Close conn and exit
			conn.Close()
			fmt.Println("Connection Closed. Bye bye.")
			os.Exit(0)
		}
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Response: %s", msg)

	}
}
