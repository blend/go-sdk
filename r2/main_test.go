package r2

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"
)

type MockTracer struct {
	StartHandler  func(*http.Request)
	FinishHandler func(*http.Request, *http.Response, time.Time, error)
}

func (mt MockTracer) Start(req *http.Request) TraceFinisher {
	if mt.StartHandler != nil {
		mt.StartHandler(req)
	}
	return MockTraceFinisher{
		Tracer: mt,
	}
}

type MockTraceFinisher struct {
	Tracer MockTracer
}

func (mtf MockTraceFinisher) Finish(req *http.Request, res *http.Response, t time.Time, err error) {
	if mtf.Tracer.FinishHandler != nil {
		mtf.Tracer.FinishHandler(req, res, t, err)
	}
}

func readString(r io.Reader) string {
	contents, _ := ioutil.ReadAll(r)
	return string(contents)
}

func mockServerOK() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK!\n")
	}))
}

var clientCert = []byte(`-----BEGIN CERTIFICATE-----
MIIDCzCCAfOgAwIBAgIQbU6bmWZG5fRA5rxB7jf9czANBgkqhkiG9w0BAQsFADAA
MCAYDzAwMDEwMTAxMDAwMDAwWhcNMjkwNTAyMDAwMDMzWjArMSkwJwYDVQQDEyAx
YzljNjViYWE0ZTU0MzdlYmE1NmQxMzRhMGU2ZWVjYzCCASIwDQYJKoZIhvcNAQEB
BQADggEPADCCAQoCggEBANd3YWl2ksmmVXlcbb/NXXiL6rjO54QOiD+Pftb/jJZv
zfazyF+rq8/4wSahh3qXm/i1yNHxmg4Iv673OizYw7viIzdMDo9TWypiPXHXJ8K0
iRBCZ8CeWcUK2I6dM2jzgmyOL/7BXtNuHpR0BVxVpFR+oiGgIrK8J8Gx8+tu9+Hj
1rHvqt8gYmgzjjwopCPfM9Z98tBnoBA3Ctxh+isRF2sumYdq+kQt1mIcyEM1Szxq
vJjX50KXQQ2qOrxH/kJQFstyNwwtCI3xIRfcAzYRo1BDkAmhUb7DSTZEkFocN8qj
EjPNAVJRKfMeSe9SQCXTGmex5jwx4+o6+xCQRJodo8sCAwEAAaNUMFIwDgYDVR0P
AQH/BAQDAgeAMBMGA1UdJQQMMAoGCCsGAQUFBwMCMCsGA1UdEQQkMCKCIDFjOWM2
NWJhYTRlNTQzN2ViYTU2ZDEzNGEwZTZlZWNjMA0GCSqGSIb3DQEBCwUAA4IBAQAl
NT2ycZ9eiTE3Y4hRTP1jizAUoJdL35Do7VZO9FLib/2G/6GgQ+XbthISvABJUbxT
S21qfv2WOnCL4WT3g3xFX3cGyQLfT+/TYXC2bbfa00vIw2YsWqpTP5/czc6hbyLh
dbsZTg4dZ0vQYQqxBPZ7v8lD6hNdaOEDBYktFlMJ+NnSWx3/bRWIKzZpnplLejnW
cX+S+eWDeAeD5t3FUUde9BYWd8ENASma8KjD0K4yt1wStmX42Evk2hPyamHM9RY1
ns/GvjpAhbx3A7aaBi3UWL06uMhJx4xEVnCbILprIwWkH9HPfWvZf//2K7h7dlM8
6BgVHMMLy9ZomrnlUee5
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIICsDCCAZigAwIBAgIRAIm5oWJbAYyx5n/E3WbVo6MwDQYJKoZIhvcNAQELBQAw
ADAgGA8wMDAxMDEwMTAwMDAwMFoXDTQ0MDUwMjAwMDAzMVowADCCASIwDQYJKoZI
hvcNAQEBBQADggEPADCCAQoCggEBAM0+B0QO1oM44pK76EssgE+fu1CbTnIS/kAG
V0NCfLTvpdym82wsJwqEnqBYsIkTpB0o2/9rHnmceIdVvOBIC1Mj0m2iG9iZyXxx
/oGHeLo0RH3hFNoI1X/zggtAgAP93OvHw7GUMmcJPFYnCi3q2wY8L0342OOxiqcY
FGNeAvuFHGFiLoRRZwh0jAg2jMdkjPADyIwr6Wt/hYjrCJPYGel7DXF/RP3wsaIC
IM97lhnPb6f1tRWrvq774Gf0FL3cbJJoqC1mfeb+5hmSmHpX9ky6FudxIvviXLtT
mSzskHa1TcV5bNV9PzQpuRajwBk3NhH5NbwpIqT3IZ/ZoTZQ3y8CAwEAAaMjMCEw
DgYDVR0PAQH/BAQDAgKkMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQAD
ggEBACJTJe29CVo3qhFMcOKvKvZr6sCbXCtHr4fUI126AWPpUiUIOO12H8VvMgAy
boouth25FK/orgbqDlbjeVPaSKDyE+smkgyzB+URSdFK0u6U1+ON/aQimdMRamme
GzxIhu3cDMWuanC61ND91hYzakgPJEDXF+RfUi7ngGn30e8QLb0EwByjD99SJnak
lyFm5+SNAkOkZZ6Vc3CalFkaCgzehzuQTLImv2SFg1W3DiFQ90A572Qn3u+cj6+7
tUcjUZrzhHC0unYqqCb7KbJvNotAXT+kYq5FDVifeH+j09mzZ/SXCkTMtko0i0+m
dg8HIlKx9z1Bxys5Ko5hGpupXDY=
-----END CERTIFICATE-----
`)

var clientKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA13dhaXaSyaZVeVxtv81deIvquM7nhA6IP49+1v+Mlm/N9rPI
X6urz/jBJqGHepeb+LXI0fGaDgi/rvc6LNjDu+IjN0wOj1NbKmI9cdcnwrSJEEJn
wJ5ZxQrYjp0zaPOCbI4v/sFe024elHQFXFWkVH6iIaAisrwnwbHz62734ePWse+q
3yBiaDOOPCikI98z1n3y0GegEDcK3GH6KxEXay6Zh2r6RC3WYhzIQzVLPGq8mNfn
QpdBDao6vEf+QlAWy3I3DC0IjfEhF9wDNhGjUEOQCaFRvsNJNkSQWhw3yqMSM80B
UlEp8x5J71JAJdMaZ7HmPDHj6jr7EJBEmh2jywIDAQABAoIBAC09+u1LIX1H+NCX
0M+iToseTfXqNACtkHxQJCD+3cVEyqmPjHZSNKxhniT/a9QY+34YpYc3xNJHkgAq
F0QNa+QKkrxssu3zYcQfhqldtRKUF+ebGe//D/ho05n2djIGV491t6w1bDTW/YLM
bce6j9vSDzciScbf7TUlqYL49QGwE9BIP+G1KijGFZMnng1O2+vUt8SOm7jzdUTH
11I2bNduDN8iUDTKPWOI4l2HtYUKQuNi+hw/xhrwS4IC3234yN10tjMvY8RJMnVC
pjTr7CzQK0ACtH6kdgcdSVl1oJqMXsVCHO76rR/Km8XL5FIsAkjMfpzs11dOeK/n
j3EoijECgYEA9OAP2TfWuASUF7qjmDZKezwaIfhZRSRwleNsuobMYvV+FPZWbaTU
NSNHOquiu6LJiv3HJUzfXkm5NLDuUQceP9lTbeCMIQCFIZ11AQQK5wYXcSjtwZ1N
SDO80gh7l4NdlYSr8rCi3lYCvgwmbSNGIOcNqbmv8bry3jMgzFBImIMCgYEA4UFJ
zIUOi/lRsdPRL6ERNcpbHe5RNdDMDWwN+HIjH/1PQRINsmwlQhzEN3g+P/bMjzSZ
2s9Edt8hz89dvht42flA0VISEbttshR0TR5vRAF22VEQ8TI98JKNF7Ge81uDJqta
1d8F7v9uVteDXoPRoYMxvnbyThUN9X5P6YgsFRkCgYBWu+pBKSsPoOeHhB8f8dLt
1Xr4H0wXVnHeVWCUrNxGDOgsqpgwW9qiO62mFVcdmOpEJeFcz96qOfi0thqjbp8D
RInteESKB/If1vKzemgWLi0tcq7MDlhqQ5EU39ZO80O5ivWQj4oQsGxmPk16CK11
SAGp5VBxkaMmmvt6AtHD2wKBgHP/YBavKneQk65kquPBKRCvPU7ji/SPqpT64RLh
DA+MLcUPm/gW0vUBxVXfWQcte9f/OX/BnrssWsgePGMK2Kg/QE7K2b1B7NJ40A9q
rdeyfVaZ9YSP3+/EOF5MPNOLe7VtJqDecbrK1TJpVyBT958Z5YL01AC7vO1/930G
f9T5AoGBALYDPM9zZDtbbh7SW/dpKpTS6BAm229gXLsMNRtmvAP3TZ4CgFuayIay
IVg93a56rlerNoUfafYkFNjy8xwNR734qtpCxn11L0L7gJtBW3/AUZIL+y3OwqtY
L8KvXHCoQKhC9uOMhb0Pju/4q3izTWrHnzPgLei+yDy1u+piDraf
-----END RSA PRIVATE KEY-----
`)
