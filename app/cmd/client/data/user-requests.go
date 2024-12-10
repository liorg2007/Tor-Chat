package main

import (
	"container/list"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"marshmello/pkg/handlers"
)

type AuthUserRequest struct {
	Username string
	Password string
}

type SendMessage struct {
	username string
	message  string
	token    string
}

type GetMessagees struct {
	token string
}

func SendRegister(nodeList *list.List, data AuthUserRequest) {
	req, err := CreateRequestThroughNetwork(nodeList, data, "auth/login")

	if err != nil {
		fmt.Println("Cant create request")
		return
	}

	respJson, err := SendHttpRequest(nodeList.Front().Value.(NodeInfo).Addr, req, "redirect")
	if err != nil {
		resp, err := DecodeErrorFromNetwork(err.Error(), nodeList)
		if err != nil {
			fmt.Println("error decoding 1")
			return
		}
		decodedResp, err := base64.StdEncoding.DecodeString(resp)

		fmt.Println(string(decodedResp))
		return
	}

	var resp handlers.EncryptedResponse
	err = json.Unmarshal(respJson, &resp)
	if err != nil {
		fmt.Println("error unmarshaloing")
		return
	}

	responseString, err := DecodeRequestThroughNetwork(nodeList, resp.Data)
	if err != nil {
		fmt.Println("error decoding 2")
		return
	}

	decodedResp, err := base64.StdEncoding.DecodeString(responseString)
	if err != nil {
		fmt.Println("error decoding b64")
		return
	}

	fmt.Println(string(decodedResp))
}
