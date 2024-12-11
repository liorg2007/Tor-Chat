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

func SendRegister(nodeList *list.List, data AuthUserRequest) error {
	req, err := CreateRequestThroughNetwork(nodeList, data, "auth/register")

	if err != nil {
		//fmt.Println("Cant create request")
		return err
	}

	respJson, err := SendHttpRequest(nodeList.Front().Value.(NodeInfo).Addr, req, "redirect")
	if err != nil {
		resp, err := DecodeErrorFromNetwork(err.Error(), nodeList)
		if err != nil {
			//fmt.Println("error decoding 1")
			return err
		}
		decodedResp, err := base64.StdEncoding.DecodeString(resp)

		if err != nil {
			return err
		}

		return fmt.Errorf(string(decodedResp))
	}

	var resp handlers.EncryptedResponse
	err = json.Unmarshal(respJson, &resp)
	if err != nil {
		//fmt.Println("error unmarshaloing")
		return err
	}

	responseString, err := DecodeRequestThroughNetwork(nodeList, resp.Data)
	if err != nil {
		//fmt.Println("error decoding 2")
		return err
	}

	_, err = base64.StdEncoding.DecodeString(responseString)
	if err != nil {
		//fmt.Println("error decoding b64")
		return err
	}

	return nil
}

func SendLogin(nodeList *list.List, data AuthUserRequest) (string, error) {
	req, err := CreateRequestThroughNetwork(nodeList, data, "auth/login")

	if err != nil {
		fmt.Println("Cant create request")
		return "", nil
	}

	respJson, err := SendHttpRequest(nodeList.Front().Value.(NodeInfo).Addr, req, "redirect")
	if err != nil {
		resp, err := DecodeErrorFromNetwork(err.Error(), nodeList)
		if err != nil {
			fmt.Println("error decoding 1")
			return "", nil
		}
		decodedResp, err := base64.StdEncoding.DecodeString(resp)

		if err != nil {
			return "", err
		}

		return "", fmt.Errorf(string(decodedResp))
	}

	var resp handlers.EncryptedResponse
	err = json.Unmarshal(respJson, &resp)
	if err != nil {
		//fmt.Println("error unmarshaloing")
		return "", nil
	}

	responseString, err := DecodeRequestThroughNetwork(nodeList, resp.Data)
	if err != nil {
		//fmt.Println("error decoding 2")
		return "", nil
	}

	decodedResp, err := base64.StdEncoding.DecodeString(responseString)
	if err != nil {
		//fmt.Println("error decoding b64")
		return "", nil
	}

	var key AuthResponse
	err = json.Unmarshal(decodedResp, &key)
	if err != nil {
		//fmt.Println("error unmarshaloing")
		return "", nil
	}

	return key.Token, nil
}

func SendMessage(nodeList *list.List, data SendMessageStruct) error {
	req, err := CreateRequestThroughNetwork(nodeList, data, "messages/send")

	if err != nil {
		//fmt.Println("Cant create request")
		return err
	}

	respJson, err := SendHttpRequest(nodeList.Front().Value.(NodeInfo).Addr, req, "redirect")
	if err != nil {
		resp, err := DecodeErrorFromNetwork(err.Error(), nodeList)
		if err != nil {
			fmt.Println("error decoding 1")
			return err
		}
		decodedResp, err := base64.StdEncoding.DecodeString(resp)

		if err != nil {
			return err
		}

		return fmt.Errorf(string(decodedResp))
	}

	var resp handlers.EncryptedResponse
	err = json.Unmarshal(respJson, &resp)
	if err != nil {
		//fmt.Println("error unmarshaloing")
		return err
	}

	responseString, err := DecodeRequestThroughNetwork(nodeList, resp.Data)
	if err != nil {
		//fmt.Println("error decoding 2")
		return err
	}

	_, err = base64.StdEncoding.DecodeString(responseString)
	if err != nil {
		//fmt.Println("error decoding b64")
		return err
	}

	return nil
}

func ReceiveMessages(nodeList *list.List, data GetMessagees) ([]MessageResponse, error) {
	req, err := CreateRequestThroughNetwork(nodeList, data, "messages/fetch")

	if err != nil {
		//fmt.Println("Cant create request")
		return nil, err
	}

	respJson, err := SendHttpRequest(nodeList.Front().Value.(NodeInfo).Addr, req, "redirect")
	if err != nil {
		resp, err := DecodeErrorFromNetwork(err.Error(), nodeList)
		if err != nil {
			//fmt.Println("error decoding 1")
			return nil, err
		}
		decodedResp, err := base64.StdEncoding.DecodeString(resp)

		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf(string(decodedResp))
	}

	var resp handlers.EncryptedResponse
	err = json.Unmarshal(respJson, &resp)
	if err != nil {
		//fmt.Println("error unmarshaloing")
		return nil, err
	}

	responseString, err := DecodeRequestThroughNetwork(nodeList, resp.Data)
	if err != nil {
		//fmt.Println("error decoding 2")
		return nil, err
	}

	decodedResp, err := base64.StdEncoding.DecodeString(responseString)
	if err != nil {
		//fmt.Println("error decoding b64")
		return nil, err
	}

	messages, err := parseMessages(string(decodedResp))
	if err != nil {
		//fmt.Println("Error:", err)
		return nil, err
	}

	return messages, nil
}
