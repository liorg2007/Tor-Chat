package torpacket

type RawMessage struct {
	Code     int
	JsonData string
}

type GetAesMsg struct {
	AesKey string
}

type RedirectMsg struct {
	Addr              string
	RedirectedMessage string
}

type ReceiveMsg struct {
	Message string
}

type AckMsg struct {
	Message string
}

type ParsedMessage struct {
	Code int         // The message type
	Data interface{} // The actual data, which will vary based on the message type
}

// Message type codes used in the Tor-like protocol.
const (
	GetAES   = 100 // Request an AES key from the node.
	Redirect = 101 // Redirect an encrypted message to another node.
	Receive  = 102 // Receive and process a message from another node.
	Destroy  = 103 // Destroy the circuit or session.
	Ack      = 104 // Acknowledge the reception of a message.
)
