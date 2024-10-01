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
