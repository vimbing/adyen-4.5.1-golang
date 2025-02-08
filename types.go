package adyen

import (
	"github.com/lestrrat-go/jwx/v2/jwk"
)

type Card struct {
	Number      string
	CVV         string
	ExpiryMonth string
	ExpiryYear  string
}

type Encryptor struct {
	Url                 string
	EncryptionKey       jwk.Key
	EncryptionKeyString string
	Card                Card
	genTime             string
}
