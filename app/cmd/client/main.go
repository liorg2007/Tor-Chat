package main

import (
	"fmt"
	"marshmello/pkg/encryption"
)

type NodeInfo struct {
	Addr            string
	AesEncryptor    encryption.AESEncryptor
	Session         string
	RedirectionAddr string
}

func main() {
	circuit, err := CreateCircuit("http://localhost:8081/", "node2:8080", "node3:8080", "192.168.99.205:1234")

	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = GetAesFromNetwork(&circuit.Circuit)

	if err != nil {
		fmt.Println(err)
		return
	}
}
