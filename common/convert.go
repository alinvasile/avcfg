package common

import (
	"bytes"
	"encoding/binary"
)

// Converts Big Endian binary format of a 4 byte integer to uint32
func toUInt32(b []byte) (uint32, error) {
	var result uint32
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.BigEndian, &result)
	return result, err
}

// Converts Big Endian binary format of a 4 byte integer to uint16
func toUInt16(b []byte) (uint16, error) {
	var result uint16
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.BigEndian, &result)
	return result, err
}

// Converts an uint32 to a 4 byte Big Endian binary format
func int32ToByteArr(i uint32) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, i)
	return buf.Bytes(), err
}

// Converts an uint16 to a 4 byte Big Endian binary format
func uint16ToByteArr(i uint16) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, i)
	return buf.Bytes(), err
}
