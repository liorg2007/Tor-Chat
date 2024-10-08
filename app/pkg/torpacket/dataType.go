package torpacket

type GetAesRequest struct {
	RsaKey string
}

type GetAesResponse struct {
	AesKey string
}

type SetRedirectRequest struct {
	Addr string
}

type RedirectRequest struct {
	Message string
}

type ReceiveRequest struct {
	Message string
}
