package torpacket

type RawMessage struct {
	code     int
	jsonData string
}

type EncryptedMessage struct {
	encryptedData string // in base64
}

type GetAesMsg struct {
	AesKey string
}
