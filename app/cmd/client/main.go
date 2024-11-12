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

	for n := nodeList.Front(); n != nil; n = n.Next() {
		node, ok := n.Value.(NodeInfo)
		if !ok {
			fmt.Println("unexpected type in node list; expected *NodeInfo")
			return
		}

		fmt.Printf("Addr: %s, Key: %s, Session: %s, Redirect: %s\n", node.Addr, node.AesEncryptor.Key, node.Session, node.RedirectionAddr)
	}
}
