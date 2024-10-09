package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"marshmello/pkg/encryption"
	"marshmello/pkg/session"
	"net/http"
)

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

	OK
*/
func SetRedirectHandler(w http.ResponseWriter, r *http.Request, sm session.SessionManager) {
	var setRedirectRequest SetRedirectRequest
	var aesDecryption encryption.AESEncryptor
	var sessionData *session.SessionData
	var err error

	err = json.NewDecoder(r.Body).Decode(&setRedirectRequest)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Return 500 Internal Server Error
		fmt.Fprintf(w, "Error reading json data.")
		return
	}

	sessionData, err = sm.PullData(setRedirectRequest.Session)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Return 500 Internal Server Error
		fmt.Fprintf(w, "Error getting session key.")
		return
	}

	aesDecryption.Key, err = base64.StdEncoding.DecodeString(sessionData.AESKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Return 500 Internal Server Error
		fmt.Fprintf(w, "Error getting aes key.")
		return
	}

	b64decodedAddr, err := aesDecryption.DecryptBase64(setRedirectRequest.Addr)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Return 500 Internal Server Error
		fmt.Fprintf(w, "Error decrypting address.")
		return
	}

	addr, err := base64.StdEncoding.DecodeString(b64decodedAddr)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Return 500 Internal Server Error
		fmt.Fprintf(w, "Error decoding address.")
		return
	}

	err = sm.AppendAddr(setRedirectRequest.Session, string(addr))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Return 500 Internal Server Error
		fmt.Fprintf(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
