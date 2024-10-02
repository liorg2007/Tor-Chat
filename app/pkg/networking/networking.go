package networking

import (
	"encoding/binary"
	"errors"
	"io"
	"marshmello/pkg/helper"
	"net"
)

func SendData(conn net.Conn, base64data string) error {
	if !helper.IsBase64Encoded(base64data) {
		return errors.New("data is not base64 encoded")
	}

	data := []byte(base64data)

	// Get the length of the data
	length := uint32(len(data))

	// Create a byte slice for the length prefix (4 bytes for uint32 in big endian)
	lengthPrefix := make([]byte, 4)

	binary.BigEndian.PutUint32(lengthPrefix, length)

	// Write the length prefix
	if _, err := conn.Write(lengthPrefix); err != nil {
		return err
	}

	// Write the actual data
	if _, err := conn.Write(data); err != nil {
		return err
	}

	return nil
}

func ReceiveData(conn net.Conn) (string, error) {
	// Read the length prefix (4 bytes)
	lengthPrefix := make([]byte, 4)
	if _, err := io.ReadFull(conn, lengthPrefix); err != nil {
		return "", err
	}

	length := binary.BigEndian.Uint32(lengthPrefix)

	// Read the data based on the length
	data := make([]byte, length)
	if _, err := io.ReadFull(conn, data); err != nil {
		return "", err
	}

	// Convert the byte slice back to a string (base64 encoded)
	base64data := string(data)

	if !helper.IsBase64Encoded(base64data) {
		return "", errors.New("data is not base64 encoded")
	}

	return base64data, nil
}
