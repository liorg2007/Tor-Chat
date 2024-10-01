package main

import (
	"fmt"
	"marshmello/pkg/encryption"
	"marshmello/pkg/torpacket"
)

func main() {
	var aes encryption.AESEncryptor
	aes.GenerateKey(32)
	fmt.Println("Key :", aes.Key)

	rawMessage, err := torpacket.SerializeGetAES(aes.Key)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Serialized Message:", rawMessage)
	}
}
