package main

import (
	"container/list"
	"encoding/base64"
	"encoding/json"
	"errors"
	"marshmello/pkg/encryption"
	"marshmello/pkg/handlers"
)

// GetAesKey requests the AES key and session token from the a node
func GetAesKey(addr string) ([]byte, string, error) {
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

func SetRedirectAddr(redirectionAddr string, nodeInfo NodeInfo) (string, error) {
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
