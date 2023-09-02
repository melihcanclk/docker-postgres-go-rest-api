package models

type TokenDetails struct {
	Token     *string
	TokenUuid string
	UserID    string
	ExpiresIn *int64
}
