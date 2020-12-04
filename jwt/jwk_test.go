package jwt

import (
	"encoding/json"
	"testing"

	"github.com/blend/go-sdk/assert"
)

const googleJWK = `{
	"alg": "RS256",
	"use": "sig",
	"kid": "0a7dc12664590c957ffaebf7b6718297b864ba91",
	"kty": "RSA",
	"e": "AQAB",
	"n": "7NfiTQcshWgrEdKbHC2e1s92kK-YX7jS3JLFIBpT8f_j_b5y3dQdtFFS4vBoVNQkwep_34x_ihYlhA3QkwaTL2XMSiedjLnubFZBUjs7G0dgGIR3F8A06Bf5KT4g2x1dKVb0Lwwqg22XIfqaS88HdU5pDwcVmq4pVMaJQgUK-xFEC_sHdfqTV8Z0uBCr9Nik_7xz68FINDYyLhehnvwph9ui-8_WeDgU_h5xrG8H7oY28y2NCtBwXxIadB-K8pHxK2srM8wTCIivdyZS80P0jZMqyxPkt4fO33-GQWvelVmR0bS4Arb3Y4bXnoAMCEao3DTm0bgeNVz39274ippJSQ"
}`

func Test_JWK_PublicKey(t *testing.T) {
	it := assert.New(t)

	var j JWK
	it.Nil(json.Unmarshal([]byte(googleJWK), &j))

	pubKey, err := j.PublicKey()
	it.Nil(err)
	it.NotNil(pubKey)
}
