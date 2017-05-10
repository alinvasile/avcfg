package common

import (
	"fmt"
	"io"
	"net"
)

const API_VERSION uint16 = 0x01

func readApiVersion(conn net.Conn) (uint16, error) {
	versionBuf := make([]byte, 2)
	_, err := conn.Read(versionBuf)
	if err != nil {
		return 0x00, err
	}
	versionData, err := toUInt16(versionBuf)

	return versionData, nil
}

func readMessageLength(conn net.Conn) (uint32, error) {
	buf := make([]byte, 4)
	_, err := conn.Read(buf)
	if err != nil {
		return 0x00, err
	}
	data, err := toUInt32(buf)

	return data, nil
}

func readMessageBody(conn net.Conn, length uint32) (string, error) {
	buf := make([]byte, length)

	reqLen := 0
	// Keep reading data from the incoming connection into the buffer until all the data promised is
	// received
	for reqLen < int(length) {
		tempreqLen, err := conn.Read(buf[reqLen:])
		reqLen += tempreqLen
		if err == io.EOF {
			return "", fmt.Errorf("Received EOF before receiving all promised data.")
		}
		if err != nil {
			return "", fmt.Errorf("Error reading: %s", err.Error())
		}
	}
	return string(buf), nil
}

func ReadMessage(conn net.Conn) (string, error) {
	versionData, err := readApiVersion(conn)

	if err != nil {
		return "", err
	}

	if versionData != 0x01 {
		return "", fmt.Errorf("Unsupported API version: %s", fmt.Sprint(versionData))
	}

	lenData, err := readMessageLength(conn)
	if err != nil {
		return "", err
	}

	readMessageBody, err := readMessageBody(conn, lenData)
	if err != nil {
		return "", fmt.Errorf("Error reading body: %s", err.Error())
	}

	return readMessageBody, nil
}

func WriteMessage(conn net.Conn, msg string) error {
	version, err := uint16ToByteArr(API_VERSION)
	_, err = conn.Write(version)
	if err != nil {
		return err
	}
	// Send the size of the message to be sent
	bytes, err := int32ToByteArr(uint32(len([]byte(msg))))
	if err != nil {
		return err
	}
	_, err = conn.Write(bytes)
	if err != nil {
		return err
	}
	// Send the message
	_, err = conn.Write([]byte(msg))
	if err != nil {
		return err
	}
	return nil
}
