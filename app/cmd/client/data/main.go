package main

import (
	"fmt"
)

func main() {
	_, err := CreateCircuit("http://localhost:8081/", "node2:8080", "node3:8080", "192.168.1.205:8080")

	if err != nil {
		fmt.Println(err)
		return
	}

	if err != nil {
		fmt.Println(err)
		return
	}
}
