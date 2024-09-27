package packet

const (
	GetAES   = 100 // Create a new circuit
	Redirect = 101
	Receive  = 102 // Custom: Redirect to another server
	Destroy  = 103 // Custom: Receive data
	Ack      = 104
)

type RawMessage struct {
	code     int
	jsonData string
}

type EncryptedMessage struct {
	length        int
	encryptedData string // in base64
}

func SerializeGetAES(aesKey string) RawMessage {

}

func SerialzieRedirect(message EncryptedMessage, addr string) RawMessage {

}

func SerializeReceive(message string) RawMessage {

}

func SerializeDestroy() RawMessage {

}

func SerializeAck() RawMessage {

}
