package main

import (
	"container/list"
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
	nodeOne, err := CreateInitialConnection("http://localhost:8081/", "node2:8080")

	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	nodeList := list.List{}
	nodeList.PushBack(nodeOne)

	err = GetAesFromNetwork(&nodeList)

	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
}
