package main

import (
	"fmt"
	"marshmello/utils/encryption"
)

func main() {
	var enc encryption.AESEncryptor
	enc.GenerateKey(32)

	msg := "Hello world"

	fmt.Println("The original message: ", msg)

	cipher, err := enc.Encrypt([]byte(msg))

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("The cipher message: ", string(cipher))

	decrypted, err := enc.Decrypt(cipher)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("The decrypted message: ", string(decrypted))
}
