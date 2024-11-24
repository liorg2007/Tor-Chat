package main

import (
	"fmt"
)

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
