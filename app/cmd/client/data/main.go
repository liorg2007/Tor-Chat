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

	var auth AuthUserRequest
	var message SendMessageStruct
	var user GetMessagees
	auth.Username = "Hello1"
	auth.Password = "1234567"

	str := SendLogin(&circuit.Circuit, auth)
	user.Token = str
	message.Token = str
	message.Username = auth.Username
	message.Message = "Hello world"
	fmt.Printf("Token: %s", str)

	SendMessage(&circuit.Circuit, message)
	SendMessage(&circuit.Circuit, message)
	SendMessage(&circuit.Circuit, message)

	ReceiveMessages(&circuit.Circuit, user)

	//_, err = GetAesFromNetwork(&circuit.Circuit)

	if err != nil {
		fmt.Println(err)
		return
	}
}
