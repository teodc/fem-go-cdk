package types

type RegisterUserPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
