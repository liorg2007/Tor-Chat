package torpacket

type GetAesRequest struct {
	RsaKey string
}

type GetAesResponse struct {
	AesKey string
}

type RedirectRequest struct {
	Addr              string
	RedirectedMessage string
}

type ReceiveRequest struct {
	Message string
}
