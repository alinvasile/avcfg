package main

import (
	"fmt"
	"github.com/alinvasile/avcfg/common"
	"github.com/patrickmn/go-cache"
	"io"
	"log"
	"net"
	"os"
)


func startTcpServer(propertyCache *cache.Cache, port string) {
	l, err := net.Listen("tcp", "0.0.0.0"+":"+port)
	if err != nil {
		fmt.Println("Error listening on tcp:", err.Error())
		os.Exit(1)
	}

	defer l.Close()
	fmt.Println("TCP Listening on " + port)
	for {

		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}


		go handleRequest(propertyCache, conn)
	}
}

// Handles incoming requests.
func handleRequest(propertyCache *cache.Cache,conn net.Conn) {
	for {
		msg, err := common.ReadMessage(conn)
		if err != nil {
			if err == io.EOF {

				conn.Close()
				return
			}
			log.Println(err)
			return
		}
		log.Printf("Message Received: %s", msg)

		value, err := readPropertyThroughCache(propertyCache, msg)


		if err != nil {
			common.WriteMessage(conn, fmt.Sprintf("Error reading: %s", err.Error()))
			return
		} else {
			log.Printf("Write response : %s", value)
			common.WriteMessage(conn, value)
		}

	}
}




