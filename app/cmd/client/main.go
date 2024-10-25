package main

import "fmt"

var CircuitInfo map[string]NodeInfo // The key is the session token

type NodeInfo struct {
	Addr            string
	AesKey          []byte
	RedirectionAddr string
}

func main() {
	key, ses, err := GetAesKey("localhost:8081")

	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	fmt.Printf("Got key: %s\n Session Token: %s\n", key, ses)
}
