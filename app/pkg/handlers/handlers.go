package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"marshmello/pkg/encryption"
	"marshmello/pkg/session"
	"net/http"
)

// Security Notes:
// 1. All data is encrypted using AES encryption
// 2. Session management is required for all endpoints except /get-aes
// 3. The AES key is exchanged securely using RSA encryption
// 4. All binary data is encoded using base64 for safe transmission

/*
General API Response format:

	{
	    "data": base64 string  // AES encrypted payload, encoded as base64
	}
*/

// EncryptResponse takes a response object or an error message, encrypts it with AES,
// and writes it to the http.ResponseWriter in the standard response format:
//
//	{
//	    "data": "base64 string"  // AES encrypted payload, encoded as base64
//	}
func EncryptResponse(w http.ResponseWriter, aesEncryptor encryption.AESEncryptor, data interface{}, statusCode int) {
	// Marshal the response data into JSON
	responseJSON, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error encoding response data.", http.StatusInternalServerError)
		return
	}

	// Encrypt the JSON response using AES
	encryptedData, err := aesEncryptor.Encrypt(responseJSON)
	if err != nil {
		http.Error(w, "Error encrypting response data.", http.StatusInternalServerError)
		return
	}

	// Encode the encrypted data into base64
	b64EncryptedData := base64.StdEncoding.EncodeToString(encryptedData)

	// Write the encrypted response in JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"data": b64EncryptedData,
	})
}

// SendResponse takes a response object or an error message and writes it as JSON
// to the http.ResponseWriter without encryption
func SendResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	// Marshal the response data into JSON
	responseJSON, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error encoding response data.", http.StatusInternalServerError)
		return
	}

	// Set the header and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(responseJSON)
}

// GetAesHandler generates and encrypts an AES key, returning it with a session token
//
// Request Type: "/get-aes"
// Request Payload (after decryption):
//
//	{
//	    "rsa_key": string  // Client's RSA public key in PEM format
//	}
//
// Response (after decryption):
//
//	{
//	    "session": string,  // Session token for subsequent requests
//	    "aes_key": string   // AES key encrypted with client's RSA public key, base64 encoded
//	}
//
// Error Responses:
// - 400 Bad Request: "Error reading json data." or "Error reading RSA key."
// - 500 Internal Server Error: "Error creating session key." or "Error encrypting AES key."
func GetAesHandler(w http.ResponseWriter, r *http.Request, sm session.SessionManager) {
	var getAesRequest GetAesRequest
	var rsaEncryptor encryption.RSAEncryptor
	var aesEncryption encryption.AESEncryptor
	var aesKey string
	var ans GetAesResponse

	// Decode the incoming JSON request
	err := json.NewDecoder(r.Body).Decode(&getAesRequest)
	if err != nil {
		http.Error(w, "Error reading json data.", http.StatusBadRequest)
		return
	}

	// Decode the RSA public key
	rsaEncryptor.PublicKey, err = encryption.DecodeRSAPublicKey(getAesRequest.RsaKey)
	if err != nil {
		http.Error(w, "Error reading RSA key.", http.StatusBadRequest)
		return
	}

	// Generate the AES key
	aesEncryption.GenerateKey()
	aesKey = encryption.EncodeAESKey(aesEncryption.Key)

	// Create session and store the AES key
	sessionToken, err := sm.CreateSession(aesKey)
	if err != nil {
		http.Error(w, "Error creating session key.", http.StatusInternalServerError)
		return
	}

	// Encrypt the AES key using the RSA public key
	encryptedKey, err := rsaEncryptor.Encrypt(aesEncryption.Key)
	if err != nil {
		http.Error(w, "Error encrypting AES key.", http.StatusInternalServerError)
		return
	}

	// Build the response with session token and encrypted AES key
	ans.Session = sessionToken
	ans.Aes_key = base64.StdEncoding.EncodeToString(encryptedKey)

	// Send the response without encryption (AES not required here)
	SendResponse(w, ans, http.StatusOK)
}

// SetRedirectHandler sets the redirect address in a session, using AES encryption for the address
//
// Request Type: "/set-redirect"
// Request Payload (after decryption):
//
//	{
//	    "session": string,  // Session token from /get-aes
//	    "addr": string      // Base64 encoded, AES encrypted redirect address (format: "ip:port")
//	}
//
// Response (after decryption):
//
//	{
//	    "message": "OK"  // Success message
//	}
//
// Error Responses:
// - 400 Bad Request: "Error reading JSON data." or "Error decrypting address."
// - 401 Unauthorized: "Error retrieving session data."
// - 500 Internal Server Error: "Error decoding AES key." or "Error decoding base64 address."
func SetRedirectHandler(w http.ResponseWriter, r *http.Request, sm session.SessionManager) {
	var setRedirectRequest SetRedirectRequest
	var aesDecryption encryption.AESEncryptor
	var sessionData *session.SessionData
	var err error

	// Decode the incoming JSON request
	err = json.NewDecoder(r.Body).Decode(&setRedirectRequest)
	if err != nil {
		SendResponse(w, map[string]string{"error": "Error reading JSON data."}, http.StatusBadRequest)
		return
	}

	// Retrieve session data, including the AES key
	sessionData, err = sm.PullData(setRedirectRequest.Session)
	if err != nil {
		SendResponse(w, map[string]string{"error": "Error retrieving session data."}, http.StatusUnauthorized)
		return
	}

	// Decode the AES key from the session
	aesDecryption.Key, err = base64.StdEncoding.DecodeString(sessionData.AESKey)
	if err != nil {
		SendResponse(w, map[string]string{"error": "Error decoding AES key."}, http.StatusInternalServerError)
		return
	}

	// Decrypt the base64-encoded address using AES
	b64decodedAddr, err := aesDecryption.DecryptBase64(setRedirectRequest.Addr)
	if err != nil {
		EncryptResponse(w, aesDecryption, map[string]string{"error": "Error decrypting address."}, http.StatusBadRequest)
		return
	}

	// Decode the base64-encoded address string (ip:port)
	addr, err := base64.StdEncoding.DecodeString(b64decodedAddr)
	if err != nil {
		EncryptResponse(w, aesDecryption, map[string]string{"error": "Error decoding base64 address."}, http.StatusInternalServerError)
		return
	}

	// Append the redirect address to the session
	err = sm.UpdateAddress(setRedirectRequest.Session, string(addr))
	if err != nil {
		EncryptResponse(w, aesDecryption, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	// Return an AES-encrypted "OK" response
	successResponse := map[string]string{
		"message": "OK",
	}
	EncryptResponse(w, aesDecryption, successResponse, http.StatusOK)
}

/*
GET /redirect

	{
	    "session": string,     // Session key for authentication
	    "data": base64 string  // AES encrypted payload, encoded as base64
	}

The decrypted 'data' field contains:

	{
	    "type": string,        // Endpoint identifier (e.g., "/get-aes", "/set-redirect")
	    "data": base64 string  // Endpoint-specific payload, left as-is
	}
*/
func RedirectHandler(w http.ResponseWriter, r *http.Request, sm session.SessionManager) {
	var redirectReq RedirectRequest
	var reqJson RedirectRequestJson

	var aesEncryptor encryption.AESEncryptor
	var err error
	var sessionData *session.SessionData

	// Decode the incoming JSON request
	err = json.NewDecoder(r.Body).Decode(&redirectReq)
	if err != nil {
		SendResponse(w, map[string]string{"error": "Error reading JSON data."}, http.StatusBadRequest)
		return
	}

	// Retrieve session data, including the AES key
	sessionData, err = sm.PullData(redirectReq.Session)
	if err != nil {
		SendResponse(w, map[string]string{"error": "Error retrieving session data."}, http.StatusUnauthorized)
		return
	}

	// Decode the AES key from the session
	aesEncryptor.Key, err = base64.StdEncoding.DecodeString(sessionData.AESKey)
	if err != nil {
		SendResponse(w, map[string]string{"error": "Error decoding AES key."}, http.StatusInternalServerError)
		return
	}

	if sessionData.Address == "" {
		SendResponse(w, map[string]string{"error": "Addr no initialzied."}, http.StatusInternalServerError)
		return
	}

	// Decrypt the base64-encoded data using AES
	b64encodedMsg, err := aesEncryptor.DecryptBase64(redirectReq.Message)
	if err != nil {
		EncryptResponse(w, aesEncryptor, map[string]string{"error": "Error decrypting data."}, http.StatusBadRequest)
		return
	}

	reqJsonString, err := base64.StdEncoding.DecodeString(b64encodedMsg)
	if err != nil {
		EncryptResponse(w, aesEncryptor, map[string]string{"error": "Error decoding b64 data."}, http.StatusInternalServerError)
		return
	}

	// Decode the incoming JSON data
	err = json.Unmarshal(reqJsonString, &reqJson)
	if err != nil {
		EncryptResponse(w, aesEncryptor, map[string]string{"error": "Error reading JSON data."}, http.StatusBadRequest)
		return
	}

	SerializeAndRedirect(w, aesEncryptor, reqJson, sessionData)
}

func SerializeAndRedirect(w http.ResponseWriter, aesEncryptor encryption.AESEncryptor, reqJson RedirectRequestJson, sessionData *session.SessionData) {
	var path string = fmt.Sprintf("http://%s/%s", sessionData.Address, reqJson.MsgType)

	// Send POST request
	resp, err := http.Post(path, "application/json", bytes.NewBuffer([]byte(reqJson.Data)))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	// Return the response
	EncryptResponse(w, aesEncryptor, resp.Body, resp.StatusCode)
}
