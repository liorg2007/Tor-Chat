package handlers

type GetAesRequest struct {
	RsaKey string
}

type GetAesResponse struct {
	Session string
	Aes_key string
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
