package handlers

type GetAesRequest struct {
	RsaKey string
}

type GetAesResponse struct {
	Session string
	Aes_key string
}

type SetRedirectRequest struct {
	Session string
	Addr    string
}

type RedirectRequest struct {
	Session string
	Message string
}

type ReceiveRequest struct {
	Session string
	Message string
}
