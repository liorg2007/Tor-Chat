package main

import (
	"fmt"
)

func main() {
	circuit, err := CreateCircuit("http://localhost:8081/", "node2:8080", "node3:8080", "192.168.1.205:8080")

	if err != nil {
		fmt.Println(err)
		return
	}

	var data AuthUserRequest
	data.Username = "Hello1"
	data.Password = "1234567"

	SendRegister(&circuit.Circuit, data)

	//_, err = GetAesFromNetwork(&circuit.Circuit)

	if err != nil {
		fmt.Println(err)
		return
	}
}
