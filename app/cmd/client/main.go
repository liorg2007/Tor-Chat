package main

import (
	"marshmello/pkg/encryption"
)

type NodeInfo struct {
	Addr            string
	AesEncryptor    encryption.AESEncryptor
	Session         string
	RedirectionAddr string
}

func main() {
	//key, ses, err := GetAesKey("localhost:8081")

	// if err != nil {
	// 	fmt.Printf("Error: %s", err)
	// 	return
	// }

	// fmt.Printf("Got key: %s\n Session Token: %s\n", key, ses)

}
