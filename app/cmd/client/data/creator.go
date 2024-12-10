package main

import (
	"container/list"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"marshmello/pkg/encryption"
	"marshmello/pkg/handlers"
)

// CreateAesRequest, creates the struct
func CreateAesRequest(encryptor *encryption.RSAEncryptor) (handlers.GetAesRequest, error) {
	var req handlers.GetAesRequest
	var err error

	// Encode the RSA public key for transmission
	req.RsaKey, err = encryption.EncodeRSAPublicKey(encryptor.PublicKey)
	if err != nil {
		return handlers.GetAesRequest{}, err
	}

	return req, nil
}

// CreateSetAddrRequest, creates the struct of CreateSetAddrRequest with the addr being encrypted

func CreateSetAddrRequest(addr string, session string, encryptor encryption.AESEncryptor) (handlers.SetRedirectRequest, error) {
	var req handlers.SetRedirectRequest
	var err error

	b64addr := base64.StdEncoding.EncodeToString([]byte(addr))

	req.Addr, err = encryptor.EncryptBase64(b64addr)

	if err != nil {
		return handlers.SetRedirectRequest{}, err
	}

	req.Session = session

	return req, nil
}

func CreateRedirectRequest(session string, redirectedJson handlers.RedirectRequestJson, encryptor encryption.AESEncryptor) (handlers.RedirectRequest, error) {
	var finalReq handlers.RedirectRequest
	var jsonString string
	var err error

	jsonBytes, err := json.Marshal(redirectedJson)

	if err != nil {
		return handlers.RedirectRequest{}, err
	}

	jsonString = base64.StdEncoding.EncodeToString(jsonBytes)

	finalReq.Message, err = encryptor.EncryptBase64(jsonString)
	if err != nil {
		return handlers.RedirectRequest{}, err
	}

	finalReq.Session = session

	return finalReq, nil
}

func CreateRequestThroughNetwork(nodeList *list.List, message interface{}, msgType string) (handlers.RedirectRequest, error) {
	var finalReq handlers.RedirectRequest
	var reqJson handlers.RedirectRequestJson

	jsonBytes, err := json.Marshal(message)

	if err != nil {
		return handlers.RedirectRequest{}, err
	}
	jsonString := base64.StdEncoding.EncodeToString(jsonBytes)

	reqJson.Data = jsonString
	reqJson.MsgType = msgType

	for n := nodeList.Back(); n != nil; n = n.Prev() {
		currentLayer, err := CreateRedirectRequest(n.Value.(NodeInfo).Session, reqJson, n.Value.(NodeInfo).AesEncryptor)
		if err != nil {
			return handlers.RedirectRequest{}, err
		}

		jsonBytes, err := json.Marshal(currentLayer)

		if err != nil {
			return handlers.RedirectRequest{}, err
		}

		jsonString := base64.StdEncoding.EncodeToString(jsonBytes)

		reqJson.Data = jsonString
		reqJson.MsgType = "redirect"

		finalReq = currentLayer
	}

	return finalReq, nil
}

func CreateCircuit(node1 string, node2 string, node3 string, finalDst string) (MessageSender, error) {
	nodeOne, err := CreateInitialConnection(node1, node2)

	if err != nil {
		fmt.Printf("Error: %s", err)
		return MessageSender{list.List{}}, err
	}

	nodeList := list.List{}
	nodeList.PushBack(nodeOne)

	newNode, err := GetAesFromNetwork(&nodeList)

	if err != nil {
		fmt.Printf("Error 2: %s", err)
		return MessageSender{list.List{}}, err
	}

	err = SetAddrFromNetwork(&nodeList, &newNode, node3)

	nodeList.PushBack(newNode)

	if err != nil {
		fmt.Printf("Error node 2: %s", err)
		return MessageSender{list.List{}}, err
	}

	newNode2, err := GetAesFromNetwork(&nodeList)

	if err != nil {
		fmt.Printf("Error node 3 setup: %s", err)
		return MessageSender{list.List{}}, err
	}

	err = SetAddrFromNetwork(&nodeList, &newNode2, finalDst)

	nodeList.PushBack(newNode2)

	if err != nil {
		fmt.Printf("Error ndoe 3 setup: %s", err)
		return MessageSender{list.List{}}, err
	}

	for n := nodeList.Front(); n != nil; n = n.Next() {
		_, ok := n.Value.(NodeInfo)
		if !ok {
			fmt.Println("unexpected type in node list; expected *NodeInfo")
			return MessageSender{list.List{}}, err
		}

		//fmt.Printf("Addr: %s, Key: %s, Session: %s, Redirect: %s\n", node.Addr, node.AesEncryptor.Key, node.Session, node.RedirectionAddr)
	}

	return MessageSender{nodeList}, nil
}
