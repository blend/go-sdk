/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import (
	"os"
	"testing"
)

// TestMain is the testing entrypoint.
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// test keys
const (
	TestTLSCert = `-----BEGIN CERTIFICATE-----
MIIC+jCCAeKgAwIBAgIRAKGQgEUjhTZMM2VMx9y92MUwDQYJKoZIhvcNAQELBQAw
EjEQMA4GA1UEChMHQWNtZSBDbzAeFw0xODAzMDkwMzIxMzhaFw0xOTAzMDkwMzIx
MzhaMBIxEDAOBgNVBAoTB0FjbWUgQ28wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw
ggEKAoIBAQDCONjExGZ+MwYZ1CosUB+sa9jS/AD0YkOi8AgiOYughLrKx5RuSsO9
ZaO0iwH987SFwAxBEiXwfLceEDgHYLGNfKQdYMCdh1yclr9yKrfpLV1SvPwT/utm
ek3ONwbJwqIrBP0dNWtfRhHhu2Gyc1JjxpqETdCUUZfuJWouVjVIxaIxLvyxYkUo
AS6SpUlUOOF3Wnre4+3x1RWRpXwns/HUFjsQBOIBo7pganxzcukTsQZWv+kJEA2o
EW33VdLQBuD59X6h1/qjx93s3AndeT5CoeVCAQ6PKXuV9z1WCpRewPpD+J89Noff
aueXIhTvxpFnB6W6VGVDQmnhEbnwA2IPAgMBAAGjSzBJMA4GA1UdDwEB/wQEAwIF
oDATBgNVHSUEDDAKBggrBgEFBQcDATAMBgNVHRMBAf8EAjAAMBQGA1UdEQQNMAuC
CWxvY2FsaG9zdDANBgkqhkiG9w0BAQsFAAOCAQEAYkkoNdditdKaEWrUjMc52QqJ
e4hbjqWT6W3bphGgYiKvnxgcDQYL3+RgEd7tGIHfgLkIiuM9efH+KJ4/jdXFWlcQ
7PoS9nGn0FwNvGdt9KCzNZSODSgQNt7FdsSpfw6Qzhn6XCwx3Bay9uF6cPap+wtX
SX6fD+az+dh0UPYoEltuKBv43+wLwsxAg18vBFuACI52NomvNw4G4uw4epBGGmp8
A0A4h9O67T/bFXchS+uIQnThZo4U/TCDu0xi/Q89xtjWff1YybwR85l85pEt1v7G
ei1eKWKYUxUU7lBMaECknLsJ4xsDKRSA5tvEDCkeQDCwTD7Msh5uGQ9itoWMlQ==
-----END CERTIFICATE-----
`
	TestTLSKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAwjjYxMRmfjMGGdQqLFAfrGvY0vwA9GJDovAIIjmLoIS6yseU
bkrDvWWjtIsB/fO0hcAMQRIl8Hy3HhA4B2CxjXykHWDAnYdcnJa/ciq36S1dUrz8
E/7rZnpNzjcGycKiKwT9HTVrX0YR4bthsnNSY8aahE3QlFGX7iVqLlY1SMWiMS78
sWJFKAEukqVJVDjhd1p63uPt8dUVkaV8J7Px1BY7EATiAaO6YGp8c3LpE7EGVr/p
CRANqBFt91XS0Abg+fV+odf6o8fd7NwJ3Xk+QqHlQgEOjyl7lfc9VgqUXsD6Q/if
PTaH32rnlyIU78aRZwelulRlQ0Jp4RG58ANiDwIDAQABAoIBAHcg8yTN6qfhmA5j
qnJ/us3BYL8Yv2UmmKHqZLLJZTFR+FjEzfBQf3s+SolE8jXYM5QOVfXbsdWuSYtx
G0y7LGzCVM+INtzo2A9cD5VxSlkF8EX9kQiaxbyXq/2eltVOQrXsW2x9BZzsl69D
hgs03QZCHSilqhgva+cwn85IJmq5bL5BMlNT1vFUgKz4QWISuBQc84PpH9R3P0oF
ur4PRJuh6Q3/GX2MF7fuNw+cweg6lNM2IlVmoH3jJo4byW+tzruv5O+/0s92CsSM
s5ywkZlgydrh1w4Irqli67y/jdDdA9zHcr+DBpVquJ1arez/ImRtKA9+FRNP4YvM
k3FOh0ECgYEA2UU+8iad7Kd7bcrhCq6AItlv51MxTp9ASDoFiCFJncTOGLdzcVNA
a+reF22XYdD32R94ldWGlIBp3MbNTyK5HYkTbwHG8414fahxg3Uy4je0NLQzHpIH
OQjaX+YFUtMDaGL7MCIDeC1FKCwfnWBRS/6xaZe3g4ne1wqZZ46DmxECgYEA5NfT
jsLSPXD5ZEz594jsOfTJ24RH4CgB69BQTd9z9AezMlTZE3fTUeXjhZRim1cs/+/4
lotnMuUEYOVRwtfJS+hqVGg1y7MFJMTo7O5RP2+SynnIrXBkZgXwKfNX0Dj4crnA
dxlUHPPFzNEZzkNMDuiwo4ERs5G+11OPD+UL6x8CgYAkaUVmQXB/44V83d4e8yWI
MZZeVwPRYEDemdKpgKKcrQm4/K19FW2baE318SjIfMO8gFiuC421P1v+YtavZ2tM
dtdp6AtWb6P8swjq9e4kGR+7IWPbwK8zMLegEKVdvv04NjZQV7LrJfMMC3D059pX
+QP0ZTec9LMCqMUSpMCLcQKBgQDGnjAnGx6AZzp9fHYECxoEX1qHpTMA8ZhhRGc+
f2/TYI9+YrgZtol57o5f1N8Utj//TxcyCoIiYTVAqCgjdUhoEque4Oe4CYOwWxtS
8LEh3sPH6pVrOz5YclT1BBi2R4wTfvb2J8yiaE3IK8A7DpvH4NvWvWJQuXGq0AI+
KG0EvwKBgB8nHRWRbNJ8admJukGb5HF2mS1tDuHi+vB1dsTydfPDyf33B1HoEG0p
mfr9uzS9ndAYCopZO33b1h65wlPP6jnIJheycn15n7HRjYezTr8cODMnJLrRotAJ
HCsYkCmGXiwJN2guZo6l/5+GqRo3SN19dZptrH/rC/wAai0+Ctqw
-----END RSA PRIVATE KEY-----`
)
