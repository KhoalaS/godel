package types

type Credentials struct {
	Expiry  int               `json:"expiry"`
	Token   string            `json:"token"`
	Headers map[string]string `json:"headers"`
}
