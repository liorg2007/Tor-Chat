package main

import (
	"container/list"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"marshmello/pkg/encryption"
	"marshmello/pkg/handlers"
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

	str, err := SetInitRedirectAddr(redirectionAddr, nodeOne)

	if err != nil {
		return NodeInfo{}, err
	}

	nodeOne.RedirectionAddr = redirectionAddr

	fmt.Printf("Response: %s", str)

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

func GetAesFromNetwork(nodeList *list.List) error {
	rsa := encryption.RSAEncryptor{}
	err := rsa.GenerateKey()

	if err != nil {
		return err
	}

	getAes, err := CreateAesRequest(&rsa)

	if err != nil {
		return err
	}

	req, err := CreateRequestThroughNetwork(nodeList, getAes, "get-aes")

	if err != nil {
		return err
	}

	respJson, err := SendHttpRequest(nodeList.Front().Value.(NodeInfo).Addr, req, "redirect")
	if err != nil {
		return err
	}

	var resp handlers.EncryptedResponse
	err = json.Unmarshal(respJson, &resp)

	if err != nil {
		return err
	}

	front := nodeList.Front()
	if front == nil {
		// Return an error if the list is empty
		return fmt.Errorf("error: nodeList is empty")
	}

	nodeInfo, ok := front.Value.(NodeInfo)
	if !ok {
		// Return an error if the value is not of type NodeInfo
		return fmt.Errorf("error: nodeList.Front().Value is not of type NodeInfo")
	}

	// Attempt decryption
	dec, err := nodeInfo.AesEncryptor.DecryptBase64(resp.Data)
	if err != nil {
		// Return the error if decryption fails
		return err
	}

	responseString, err := base64.StdEncoding.DecodeString(dec)

	if err != nil {
		return err
	}

	fmt.Printf("Response: %s", responseString)

	//(nodeList.Back().Value.(NodeInfo)).AesEncryptor.Key =

	return nil
}
