package models

type Data struct {
	IV         string `json:"iv" validate:"required,notblank,base64"`
	Cyphertext string `validate:"required,notblank,base64"`
}

type Response struct {
	Plaintext string `json:"plaintext"`
}
