package handlers

type GetAesRequest struct {
	RsaKey string
}

type AuthUserRequest struct {
	Username string
	Password string
}

type SendMessage struct {
	Username string
	Message  string
	Token    string
}

type GetMessages struct {
	token string
}

type GetAesResponse struct {
	Session string
	Aes_key string
}

type RegularResponse struct {
	Message string
}

type EncryptedResponse struct {
	Data string
}

type SetRedirectRequest struct {
	Session string
	Addr    string
}

type RedirectRequest struct {
	Session string
	Message string //base64
}

type RedirectRequestJson struct {
	MsgType string
	Data    string
}

type ReceiveRequest struct {
	Session string
	Message string
}
