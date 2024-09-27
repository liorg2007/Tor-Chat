package packet

const (
	GetAES   = 100 // Create a new circuit
	Destroy  = 101
	Redirect = 102 // Custom: Redirect to another server
	Receive  = 103 // Custom: Receive data
)

type RawMessage struct {
	code     int
	jsonData string
}

type EncryptedMessage struct {
	length        int
	encryptedData string // in base64
}
