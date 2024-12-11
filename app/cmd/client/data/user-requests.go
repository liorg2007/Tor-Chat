package main

import (
	"container/list"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"marshmello/pkg/handlers"
	"time"
)

type AuthUserRequest struct {
	Username string
	Password string
}

type SendMessageStruct struct {
	Username string
	Message  string
	Token    string
}

type GetMessagees struct {
	Token string
}

type AuthResponse struct {
	Token string
}

type MessageResponse struct {
	Username   *string     `json:"username"`
	Message    *string     `json:"message"`
	CreateTime *CustomTime `json:"createdAt"`
}

type MessagesContainer struct {
	Messages []MessageResponse `json:"messages"`
}

type CustomTime struct {
	time.Time
}

func (t *CustomTime) UnmarshalJSON(data []byte) error {
	// Remove quotes from the JSON string
	s := string(data[1 : len(data)-1])

	// Try parsing with microseconds
	parsedTime, err := time.Parse("2006-01-02T15:04:05.000000", s)
	if err != nil {
		return err
	}

	*t = CustomTime{parsedTime}
	return nil
}

func parseMessages(jsonData string) ([]MessageResponse, error) {
	// Create a container to unmarshal the JSON
	var container MessagesContainer

	// Unmarshal the JSON data
	err := json.Unmarshal([]byte(jsonData), &container)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	// Return the list of messages
	return container.Messages, nil
}

func SendRegister(nodeList *list.List, data AuthUserRequest) {
	req, err := CreateRequestThroughNetwork(nodeList, data, "auth/register")

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

func SendLogin(nodeList *list.List, data AuthUserRequest) string {
	req, err := CreateRequestThroughNetwork(nodeList, data, "auth/login")

	if err != nil {
		fmt.Println("Cant create request")
		return ""
	}

	respJson, err := SendHttpRequest(nodeList.Front().Value.(NodeInfo).Addr, req, "redirect")
	if err != nil {
		resp, err := DecodeErrorFromNetwork(err.Error(), nodeList)
		if err != nil {
			fmt.Println("error decoding 1")
			return ""
		}
		decodedResp, err := base64.StdEncoding.DecodeString(resp)

		fmt.Println(string(decodedResp))
		return ""
	}

	var resp handlers.EncryptedResponse
	err = json.Unmarshal(respJson, &resp)
	if err != nil {
		fmt.Println("error unmarshaloing")
		return ""
	}

	responseString, err := DecodeRequestThroughNetwork(nodeList, resp.Data)
	if err != nil {
		fmt.Println("error decoding 2")
		return ""
	}

	decodedResp, err := base64.StdEncoding.DecodeString(responseString)
	if err != nil {
		fmt.Println("error decoding b64")
		return ""
	}

	var key AuthResponse
	err = json.Unmarshal(decodedResp, &key)
	if err != nil {
		fmt.Println("error unmarshaloing")
		return ""
	}

	return key.Token
}

func SendMessage(nodeList *list.List, data SendMessageStruct) {
	req, err := CreateRequestThroughNetwork(nodeList, data, "messages/send")

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

func ReceiveMessages(nodeList *list.List, data GetMessagees) {
	req, err := CreateRequestThroughNetwork(nodeList, data, "messages/fetch")

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

	messages, err := parseMessages(string(decodedResp))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print out the parsed messages
	for _, msg := range messages {
		if msg.Username != nil {
			fmt.Printf("Username: %s\n", *msg.Username)
		}
		if msg.Message != nil {
			fmt.Printf("Message: %s\n", *msg.Message)
		}
		if msg.CreateTime != nil {
			// Format the time in a more readable way
			fmt.Printf("Created At: %s\n", msg.CreateTime.Format("2006-01-02 15:04:05"))
		}
	}
}
