package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"marshmello/pkg/encryption"
	"marshmello/pkg/session"
	"net/http"
)

// EncryptResponse takes a response object or an error message, encrypts it with AES, and writes it to the http.ResponseWriter
func EncryptResponse(w http.ResponseWriter, aesEncryptor encryption.AESEncryptor, data interface{}) {
	// Marshal the response data into JSON
	responseJSON, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := map[string]string{"error": "Error encoding response data."}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Encrypt the JSON response using AES
	encryptedData, err := aesEncryptor.Encrypt(responseJSON)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := map[string]string{"error": "Error encrypting response data."}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Encode the encrypted data into base64
	b64EncryptedData := base64.StdEncoding.EncodeToString(encryptedData)

	// Write the encrypted response in JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"data": b64EncryptedData,
	})
}

/*
Request
GET /get-aes
json:

	{
		rsa_key: b64 rsa key
	}

Response

	json:
	{
		session: string session key
		aes_key: b64 encrypted key
	}
*/
func GetAesHandler(w http.ResponseWriter, r *http.Request, sm session.SessionManager) {
	var getAesRequest GetAesRequest
	var rsaEncryptor encryption.RSAEncryptor
	var aesEncyption encryption.AESEncryptor
	var aes_key string
	var ans GetAesResponse

	err := json.NewDecoder(r.Body).Decode(&getAesRequest)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Return 500 Internal Server Error
		fmt.Fprintf(w, "Error reading json data.")
		return
	}

	rsaEncryptor.PublicKey, err = encryption.DecodeRSAPublicKey(getAesRequest.RsaKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Return 500 Internal Server Error
		fmt.Fprintf(w, "Error reading rsa key.")
		return
	}

	aesEncyption.GenerateKey()
	aes_key = encryption.EncodeAESKey(aesEncyption.Key)

	sessionToken, err := sm.CreateSession(aes_key)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Return 500 Internal Server Error
		fmt.Fprintf(w, "Error creating session key.")
		return
	}

	encryptedKey, err := rsaEncryptor.Encrypt(aesEncyption.Key)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Return 500 Internal Server Error
		fmt.Fprintf(w, "Error creating encrypted key.")
		return
	}

	ans.Session = sessionToken
	ans.Aes_key = base64.StdEncoding.EncodeToString(encryptedKey)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(ans)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

/*
Request
GET /set-redirect
json:

	{
		session: string session key
		addr: base64 string, the string is encrypted with aes.
		The decrypted string is a base64 of the addr of redirect node.
		The format is ip:port
	}

Response

	{
		data: base64 encrypted "OK"
	}
*/
func SetRedirectHandler(w http.ResponseWriter, r *http.Request, sm session.SessionManager) {
	var setRedirectRequest SetRedirectRequest
	var aesDecryption encryption.AESEncryptor
	var sessionData *session.SessionData
	var err error

	// Parse incoming request JSON
	err = json.NewDecoder(r.Body).Decode(&setRedirectRequest)
	if err != nil {
		EncryptResponse(w, aesDecryption, map[string]string{"error": "Error reading json data."})
		return
	}

	// Get the session data (AES key and other details)
	sessionData, err = sm.PullData(setRedirectRequest.Session)
	if err != nil {
		EncryptResponse(w, aesDecryption, map[string]string{"error": "Error retrieving session data."})
		return
	}

	// Decode AES key from the session
	aesDecryption.Key, err = base64.StdEncoding.DecodeString(sessionData.AESKey)
	if err != nil {
		EncryptResponse(w, aesDecryption, map[string]string{"error": "Error decoding AES key."})
		return
	}

	// Decrypt the incoming addr using AES key
	b64decodedAddr, err := aesDecryption.DecryptBase64(setRedirectRequest.Addr)
	if err != nil {
		EncryptResponse(w, aesDecryption, map[string]string{"error": "Error decrypting address."})
		return
	}

	// Decode the base64-encoded address (ip:port)
	addr, err := base64.StdEncoding.DecodeString(b64decodedAddr)
	if err != nil {
		EncryptResponse(w, aesDecryption, map[string]string{"error": "Error decoding base64 address."})
		return
	}

	// Append the redirect address to the session
	err = sm.AppendAddr(setRedirectRequest.Session, string(addr))
	if err != nil {
		EncryptResponse(w, aesDecryption, map[string]string{"error": err.Error()})
		return
	}

	// Return an AES-encrypted "OK" response
	successResponse := map[string]string{
		"message": "OK",
	}

	EncryptResponse(w, aesDecryption, successResponse)
}
