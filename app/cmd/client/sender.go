package main

import (
	"container/list"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"marshmello/pkg/encryption"
	"marshmello/pkg/handlers"
	"strings"
)

func CreateInitialConnection(addr string, redirectionAddr string) (NodeInfo, error) {
	key, ses, err := GetInitAesKey("localhost:8081")

	if err != nil {
		fmt.Printf("Error: %s", err)
		return NodeInfo{}, err
	}

	fmt.Printf("Got key: %s\n Session Token: %s\n Key Size: %d\n", key, ses, len(key))

	enc := encryption.AESEncryptor{Key: key}

	nodeOne := NodeInfo{
		Addr:         "localhost:8081",
		AesEncryptor: enc,
		Session:      ses,
	}

	_, err = SetInitRedirectAddr(redirectionAddr, nodeOne)

	if err != nil {
		return NodeInfo{}, err
	}

	nodeOne.RedirectionAddr = redirectionAddr

	return nodeOne, nil
}

// GetAesKey requests the AES key and session token from the a node
func GetInitAesKey(addr string) ([]byte, string, error) {
	var req handlers.GetAesRequest
	var rsaEnc encryption.RSAEncryptor
	var res handlers.GetAesResponse

	// Generate RSA keys
	err := rsaEnc.GenerateKey()
	if err != nil {
		return nil, "", err
	}

	// Encode the RSA public key for transmission
	req.RsaKey, err = encryption.EncodeRSAPublicKey(rsaEnc.PublicKey)
	if err != nil {
		return nil, "", err
	}

	// Send request to /get-aes
	respData, err := SendHttpRequest(addr, req, "get-aes")
	if err != nil {
		return nil, "", err
	}

	// Unmarshal the response
	if err := json.Unmarshal(respData, &res); err != nil {
		return nil, "", errors.New("error decoding AES response")
	}

	// Decode the base64-encoded AES key from the response
	aesKey, err := base64.StdEncoding.DecodeString(res.Aes_key)
	if err != nil {
		return nil, "", err
	}

	decrypted, err := rsaEnc.Decrypt(aesKey)
	if err != nil {
		return nil, "", err
	}

	return decrypted, res.Session, nil
}

func SetInitRedirectAddr(redirectionAddr string, nodeInfo NodeInfo) (string, error) {
	req, err := CreateSetAddrRequest(redirectionAddr, nodeInfo.Session, nodeInfo.AesEncryptor)

	if err != nil {
		return "", err
	}

	// Send request to /get-aes
	respJson, err := SendHttpRequest(nodeInfo.Addr, req, "set-redirect")
	if err != nil {
		return "", err
	}

	var resp handlers.EncryptedResponse
	err = json.Unmarshal(respJson, &resp)

	if err != nil {
		return "", err
	}

	dec, err := nodeInfo.AesEncryptor.DecryptBase64(resp.Data)

	if err != nil {
		return "", err
	}

	responseString, err := base64.StdEncoding.DecodeString(dec)

	if err != nil {
		return "", err
	}

	return string(responseString), nil
}

func GetAesFromNetwork(nodeList *list.List) (NodeInfo, error) {
	var res handlers.GetAesResponse
	rsa := encryption.RSAEncryptor{}
	err := rsa.GenerateKey()
	if err != nil {
		return NodeInfo{}, err
	}

	getAes, err := CreateAesRequest(&rsa)
	if err != nil {
		return NodeInfo{}, err
	}

	req, err := CreateRequestThroughNetwork(nodeList, getAes, "get-aes")
	if err != nil {
		return NodeInfo{}, err
	}

	respJson, err := SendHttpRequest(nodeList.Front().Value.(NodeInfo).Addr, req, "redirect")
	if err != nil {
		resp, err := DecodeErrorFromNetwork(err.Error(), nodeList)
		if err != nil {
			return NodeInfo{}, err
		}
		return NodeInfo{}, fmt.Errorf("error: %s", resp)
	}

	var resp handlers.EncryptedResponse
	err = json.Unmarshal(respJson, &resp)
	if err != nil {
		return NodeInfo{}, err
	}

	responseString, err := DecodeRequestThroughNetwork(nodeList, resp.Data)
	if err != nil {
		return NodeInfo{}, err
	}

	decodedResp, err := base64.StdEncoding.DecodeString(responseString)
	if err != nil {
		return NodeInfo{}, err
	}

	// Unmarshal the response
	if err := json.Unmarshal(decodedResp, &res); err != nil {
		return NodeInfo{}, errors.New("error decoding AES response")
	}

	// Decode the base64-encoded AES key from the response
	aesKey, err := base64.StdEncoding.DecodeString(res.Aes_key)
	if err != nil {
		return NodeInfo{}, err
	}

	decrypted, err := rsa.Decrypt(aesKey)
	if err != nil {
		return NodeInfo{}, err
	}

	back, ok := nodeList.Back().Value.(NodeInfo)
	if !ok {
		return NodeInfo{}, errors.New("unexpected type in node list; expected *NodeInfo")
	}

	// Create a new NodeInfo entity with the decrypted AES key and session
	newNode := NodeInfo{
		AesEncryptor: encryption.AESEncryptor{Key: decrypted},
		Session:      res.Session,
		Addr:         back.RedirectionAddr,
	}

	return newNode, nil
}

func SetAddrFromNetwork(nodeList *list.List, newNode *NodeInfo, redirectionAddr string) error {
	setAddrReq, err := CreateSetAddrRequest(redirectionAddr, newNode.Session, newNode.AesEncryptor)

	if err != nil {
		return err
	}

	req, err := CreateRequestThroughNetwork(nodeList, setAddrReq, "set-redirect")
	if err != nil {
		return err
	}

	respJson, err := SendHttpRequest(nodeList.Front().Value.(NodeInfo).Addr, req, "redirect")
	if err != nil {
		resp, err := DecodeErrorFromNetwork(err.Error(), nodeList)
		if err != nil {
			return err
		}
		return fmt.Errorf("error: %s", resp)
	}

	var resp handlers.EncryptedResponse
	err = json.Unmarshal(respJson, &resp)
	if err != nil {
		return err
	}

	responseStringEnc, err := DecodeRequestThroughNetwork(nodeList, resp.Data)
	if err != nil {
		return err
	}

	dec, err := newNode.AesEncryptor.DecryptBase64(responseStringEnc)

	if err != nil {
		return err
	}

	responseString, err := base64.StdEncoding.DecodeString(dec)

	if err != nil {
		return err
	}

	log.Println(string(responseString))

	newNode.RedirectionAddr = redirectionAddr

	return nil
}

func DecodeRequestThroughNetwork(nodeList *list.List, response string) (string, error) {
	var err error
	data := response

	for n := nodeList.Front(); n != nil; n = n.Next() {
		nodeInfo, ok := n.Value.(NodeInfo)
		if !ok {
			return "", fmt.Errorf("error: nodeList.Front().Value is not of type NodeInfo")
		}

		// Decrypt the data using AesEncryptor
		data, err = nodeInfo.AesEncryptor.DecryptBase64(data)
		if err != nil {
			return "", err
		}

		// Decode the decrypted data from Base64
		decodedData, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			return "", fmt.Errorf("error decoding Base64: %v", err)
		}

		// Unmarshal the decoded JSON data into EncryptedResponse struct
		var encryptedResponse handlers.EncryptedResponse
		if err := json.Unmarshal(decodedData, &encryptedResponse); err != nil {
			return data, nil
		}

		if encryptedResponse.Data == "" {
			return data, nil
		}

		// Use the Data field for the next iteration
		data = encryptedResponse.Data
	}

	return data, nil
}

func DecodeErrorFromNetwork(errorStr string, nodeList *list.List) (string, error) {
	// Split the string to extract the JSON part
	jsonPart := strings.TrimPrefix(errorStr, "HTTP error: ")

	// Define a struct to hold the parsed data
	var result struct {
		Data string `json:"Data"`
	}

	// Parse the JSON
	err := json.Unmarshal([]byte(jsonPart), &result)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return "", err
	}

	resp, err := DecodeRequestThroughNetwork(nodeList, result.Data)
	if err != nil {
		return "", err
	}
	return resp, nil
}
