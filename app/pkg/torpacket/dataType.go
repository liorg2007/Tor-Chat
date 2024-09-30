package torpacket

type RawMessage struct {
	code     int
	jsonData string
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
