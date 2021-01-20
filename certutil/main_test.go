package certutil

var (
	caCertLiteral = []byte(`-----BEGIN CERTIFICATE-----
MIIDgTCCAmmgAwIBAgIQUVFYoBt9rT3UbtDB2PPYsjANBgkqhkiG9w0BAQsFADBq
MRYwFAYDVQQGEw1Vbml0ZWQgU3RhdGVzMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYw
FAYDVQQHEw1TYW4gRnJhbmNpc2NvMQ8wDQYDVQQKEwZXYXJkZW4xEjAQBgNVBAMT
CXdhcmRlbi1jYTAeFw0xODEwMTUyMzIwMjlaFw0yODEwMTUyMzIwMjlaMGoxFjAU
BgNVBAYTDVVuaXRlZCBTdGF0ZXMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNV
BAcTDVNhbiBGcmFuY2lzY28xDzANBgNVBAoTBldhcmRlbjESMBAGA1UEAxMJd2Fy
ZGVuLWNhMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAocLIhZrHbvjZ
dyJSaJilR5QgTw7ejS8aM8VKbidBrfOMXKTgm0mRvP6hknK34ychGke7Le9BiD5a
eGrWXn4TMxrBkbygo/yVjhSqMfBKFn0xBIBLGneynoWaygUKn/Cbz8y9J4NJ4yVO
sWMX2VGGHH02tIamSD1TtzK/Q0zpYnZfrubG0UDkC1u8W6+vaMkLBGCCSMqvclE4
SJONtcTLBG9Y5tU4wUNqbY3JSy1vjySQVk5XEPc8tyjyJiBJsSDIMRd57CQywRKs
xUUsThQ5IuyhJFn1UojNMr0BM4P5wKN5V+7AsDEaBTV9IKziY3LOcvPT334pt3pv
plPvxPwyNwIDAQABoyMwITAOBgNVHQ8BAf8EBAMCAqQwDwYDVR0TAQH/BAUwAwEB
/zANBgkqhkiG9w0BAQsFAAOCAQEAVJjAVBmO5PqhM5fQVA2lGQaFuWE0DraoFKAC
BDLCLvUmqtSKgjmpgnMUPmGaVWvzjD02lOEGv2h59iwNYG1WUqqTJEDdQHw9ucxP
FsYXt00UKB9YbEtiK+/97+rYU8EabMB/17q9ckoAoQ38cwVdKOKL0E2wRC86sfCn
sr4omD+d7t6NAv7YsIUmyxtnzE0iS1lEXQfZpOcPpIx1K554YK9vX+DJS6lBsxMS
pST+Wru/8PRGk6jOEnmFwHEut6JCORk8JqBcd7KXobJ4G3qPEFcr2Fu/ogrPj4tN
0+JpKvBNKgdklRFSZ6O3BudGXjBrMMQno8z8pcSlKRyKYLl2dA==
-----END CERTIFICATE-----`)
	caKeyLiteral = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAocLIhZrHbvjZdyJSaJilR5QgTw7ejS8aM8VKbidBrfOMXKTg
m0mRvP6hknK34ychGke7Le9BiD5aeGrWXn4TMxrBkbygo/yVjhSqMfBKFn0xBIBL
GneynoWaygUKn/Cbz8y9J4NJ4yVOsWMX2VGGHH02tIamSD1TtzK/Q0zpYnZfrubG
0UDkC1u8W6+vaMkLBGCCSMqvclE4SJONtcTLBG9Y5tU4wUNqbY3JSy1vjySQVk5X
EPc8tyjyJiBJsSDIMRd57CQywRKsxUUsThQ5IuyhJFn1UojNMr0BM4P5wKN5V+7A
sDEaBTV9IKziY3LOcvPT334pt3pvplPvxPwyNwIDAQABAoIBAE44jZSluui7pxUH
4Me5gLrzMwhnzVCCRW5hadOIYidGzUoDZKVS1wmILANQPAG46zX1F9tnhc3nMQJn
z0+UMZ6VwqLBjwCdZIh3D7HkWz8yvKBLlYWv+gzT2Hz2ISQkSoQaoiFYSENxx/LZ
UzAATxS5q9lx/ziRh0jgXkKiD31PsWU5GRBOlrq8PftrbxmXsej5Tnw2arJubZLD
Rv6jDbV86S5NHCmOP8qswcsYyl2FY7GB2Zwv+02iLOFZMFugUDFqf+z2IiUhBYr2
0vAl1SP6aEEDaarLCksNL0j+60JfoJw8pJFsdYi4rEXaOogfb4ujOyRjJ7DaE+L7
JseO4ekCgYEA14cvHh5IQFIVuGCucrwl60+RvhwiXmWNf1kDkYSDIj1anU7SV84U
OktIL9To14ARyEBddCIsg9p9vMhMwdONcviB0Scl965MIKxdSWbtO94VkguYmOSY
MMncs+h2uFEiqM4BWtPvFx7sm+ym40mfwf7x9hB3rL7ILCXJ9NK18x0CgYEAwCLq
Q+5KdAliHwdDjZ+yqXBcayBdEht3Kske1JYuwxVV7trnIIS47muAq4zAnKDmQxDp
fG5CC+H68TltHsOdj6Sx6+LK2IDgnhyBcR/erDc7d3gcWrDmnBX3v+kWk1gaJLA9
QCK0thOHJR5rNfSshv7xUxBhzzKJZjBGkCx2hmMCgYEAw8Q+wBSxe/sTT45B5mWP
69UyhIP5k9SaWfn4i8zZb9ha3lgiJy9AoFKRFyFE/bmObV5JhJsl4/4qB7fIQKZv
7OZcxCXTDs56x5LIiUu8YIyw+x8dVIMO2gIBPkkAzRqlaL717BJlMZMdR+QFEp5l
RkjUbrU2fuor2C3a604ZhuECgYAOd9KaMOxsVnSMD5j3pinm9m4PJw27GlRb5W8r
5O86g3XIGbXDzhq64V98C4pJgRg1vqVGWGsix+6EkaV05pgLxAQd1iMeMH45ib35
jcPPxgz1TxcbgSysXj8ctWmcyMqXLqo+FR70gv3vPp4mHvmK5NG5RoeG+bA9AK+f
lamLLwKBgHle/24Q+I95j0lqamYwlZhILUrGTxCKvcEoX6xy6OMNjr1/l3EoQWf0
MDk9AlYLsDFm3F/9AUZY4F4KdWFjvj3emZ+cfw048+U2W3UM6z9G6Gw/ro69CRxD
zK+1ATHUx9OfrF80qQmjFbVUBTmUi0qm4i5iztO+/sICG9P5rFr5
-----END RSA PRIVATE KEY-----`)

	certLiteralCommonName = "56cb6b4422fc4215b2d0696adde0a74a"
	certLiteral           = []byte(`-----BEGIN CERTIFICATE-----
MIIDnjCCAoagAwIBAgIRANSKBwYVW3egUP3+MP2bxI4wDQYJKoZIhvcNAQELBQAw
ajEWMBQGA1UEBhMNVW5pdGVkIFN0YXRlczETMBEGA1UECBMKQ2FsaWZvcm5pYTEW
MBQGA1UEBxMNU2FuIEZyYW5jaXNjbzEPMA0GA1UEChMGV2FyZGVuMRIwEAYDVQQD
Ewl3YXJkZW4tY2EwHhcNMTgxMDE1MjMyMDI5WhcNMTkxMDE1MjMyMDI5WjCBgTEW
MBQGA1UEBhMNVW5pdGVkIFN0YXRlczETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQG
A1UEBxMNU2FuIEZyYW5jaXNjbzEPMA0GA1UEChMGV2FyZGVuMSkwJwYDVQQDEyA1
NmNiNmI0NDIyZmM0MjE1YjJkMDY5NmFkZGUwYTc0YTCCASIwDQYJKoZIhvcNAQEB
BQADggEPADCCAQoCggEBAMGv7R81NsvUkFhvpFOmsG60ieFE5bBLABCku9FB0BGy
XzAemt75Mqd2xuXTqbKGaPxwz5xIVXmwmbaT1hd7HN/Laq4dckA/PaalVrrbVoxs
7Fczd2T99ui4oB8ROyAbb3FS0oy8JBye3GGAu02WpqwVFhp4PZ3e1z2ecEqBvDcW
p7Vt7zPvu6JHtxV7tIO9lkIkjJCzR7Aj8YBcTnv1fZ7NkqDKfz0B8RTQe2SAmEdS
SEb3ivdj/5l0f615MNtW44Y5XKYm9S4te6HfugJLw1H9DKmA4FtpzZt9ZYSKYFHz
Thx3NDEcIvgSdnhnwbTsX/coGxvLat12eMLHjemKI2ECAwEAAaMnMCUwDgYDVR0P
AQH/BAQDAgeAMBMGA1UdJQQMMAoGCCsGAQUFBwMCMA0GCSqGSIb3DQEBCwUAA4IB
AQCEJCFRtIFwwanpayhzlawNyRG5z+nJNhltWYywylIMjtfJrn0/dkXCbUwm5gU3
WGVOvnAnUvEBRnsbArZHiEMr6W1SAH4gRgCIy61kHjNbyJYHL6sVg/6qdefpO2bv
ZMdF25yG1R5AJfmjhWOaRBkAea+hp1eKviRVJORZD46lQA/loQXFW/eA6fKK7qHU
UDoUR8Yx5F35wYo0D5X9rHVIhzwTlYuVm7pjwA6yGMt2LMlk+kL33UAa0nM4U7EP
yhH7Giu+hrm4cLTg08myUCUsHfY9InX0d/XdWe3VAymvEcjATEndb59aCDvezSql
J7m6mAXNucJmI64JC6YGFx9P
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIDgTCCAmmgAwIBAgIQUVFYoBt9rT3UbtDB2PPYsjANBgkqhkiG9w0BAQsFADBq
MRYwFAYDVQQGEw1Vbml0ZWQgU3RhdGVzMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYw
FAYDVQQHEw1TYW4gRnJhbmNpc2NvMQ8wDQYDVQQKEwZXYXJkZW4xEjAQBgNVBAMT
CXdhcmRlbi1jYTAeFw0xODEwMTUyMzIwMjlaFw0yODEwMTUyMzIwMjlaMGoxFjAU
BgNVBAYTDVVuaXRlZCBTdGF0ZXMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNV
BAcTDVNhbiBGcmFuY2lzY28xDzANBgNVBAoTBldhcmRlbjESMBAGA1UEAxMJd2Fy
ZGVuLWNhMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAocLIhZrHbvjZ
dyJSaJilR5QgTw7ejS8aM8VKbidBrfOMXKTgm0mRvP6hknK34ychGke7Le9BiD5a
eGrWXn4TMxrBkbygo/yVjhSqMfBKFn0xBIBLGneynoWaygUKn/Cbz8y9J4NJ4yVO
sWMX2VGGHH02tIamSD1TtzK/Q0zpYnZfrubG0UDkC1u8W6+vaMkLBGCCSMqvclE4
SJONtcTLBG9Y5tU4wUNqbY3JSy1vjySQVk5XEPc8tyjyJiBJsSDIMRd57CQywRKs
xUUsThQ5IuyhJFn1UojNMr0BM4P5wKN5V+7AsDEaBTV9IKziY3LOcvPT334pt3pv
plPvxPwyNwIDAQABoyMwITAOBgNVHQ8BAf8EBAMCAqQwDwYDVR0TAQH/BAUwAwEB
/zANBgkqhkiG9w0BAQsFAAOCAQEAVJjAVBmO5PqhM5fQVA2lGQaFuWE0DraoFKAC
BDLCLvUmqtSKgjmpgnMUPmGaVWvzjD02lOEGv2h59iwNYG1WUqqTJEDdQHw9ucxP
FsYXt00UKB9YbEtiK+/97+rYU8EabMB/17q9ckoAoQ38cwVdKOKL0E2wRC86sfCn
sr4omD+d7t6NAv7YsIUmyxtnzE0iS1lEXQfZpOcPpIx1K554YK9vX+DJS6lBsxMS
pST+Wru/8PRGk6jOEnmFwHEut6JCORk8JqBcd7KXobJ4G3qPEFcr2Fu/ogrPj4tN
0+JpKvBNKgdklRFSZ6O3BudGXjBrMMQno8z8pcSlKRyKYLl2dA==
-----END CERTIFICATE-----`)

	keyLiteral = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAwa/tHzU2y9SQWG+kU6awbrSJ4UTlsEsAEKS70UHQEbJfMB6a
3vkyp3bG5dOpsoZo/HDPnEhVebCZtpPWF3sc38tqrh1yQD89pqVWuttWjGzsVzN3
ZP326LigHxE7IBtvcVLSjLwkHJ7cYYC7TZamrBUWGng9nd7XPZ5wSoG8NxantW3v
M++7oke3FXu0g72WQiSMkLNHsCPxgFxOe/V9ns2SoMp/PQHxFNB7ZICYR1JIRveK
92P/mXR/rXkw21bjhjlcpib1Li17od+6AkvDUf0MqYDgW2nNm31lhIpgUfNOHHc0
MRwi+BJ2eGfBtOxf9ygbG8tq3XZ4wseN6YojYQIDAQABAoIBABiP7dsyTWOl7jQ4
3Db7gY5YeM/Hg/VKXZS+v063MOK9oxKgHvW91m2kQ27r265XG0NALyPbjHNlMOkV
cGYD59J9oma4Nz/shS3387q4jA481e/tB+wXxNMYbr3h2oSk1goh/a95QH8cqkf3
IkmhnDtgZTAwJWg61ULsL8NTFoJJtUNJR4Ujab7WfHK+nZ4WH969M8GmZPps0TT5
uCRplZXoufq5yc6P38pyXN6iMDAl+Fv6mHQ+7vNGFUaTF9hru60nf+Rk5QN3rY6Q
Ht3rVoFs06DAn+RkGxtjoKRgFXadT6182nOCbUkDpUHOrzcgwaKUnFJbTCvV/UMx
nz09rU0CgYEA/TKfLt4iYmAswKx4yElJtI6lpeAB/RVr8WvEXa5FPXdOJFwvd4qR
16xlmgVoRrj/Q1clJrfDCuNwQUp/f4+iZzNkfXOXKVu544EqQKlMhsF54SIWIDWl
AJAEGcrnvIZp+c7goPUVGxkb7ByJyaexXeAF42jIV/QKgnu+022dcFMCgYEAw9Sx
5TFlRTqVFfnIDEWsOImp4uAP7aXYQw00A3c78LMGEnZjwP6PkY8KBwMlyOfNAzaj
Z4S822KKZxZbqZwsnNLfT//tywKTGsED6jT7hvXACMQPN2I79e9zoI4mipO6qn6r
zMsUCSasa9HP1HMgT0VsHawP9sEQwKzm1RImtvsCgYAG6148SqfH5nbyoQP//Ti9
bXSLbu6++tnjVB2erceIoX0KM1a8vpWzvitcpS8vV5jqPzRttXHoF0UXE7EeTE+U
4GurngAQikgvNnVwJPBJcnohM+xE3xJuTIFALzJcDZRR1lx/KZN1FR+VOUZus11+
gkPG7jMjbDDpVfZmEsXNrwKBgQC9eaQOL/bePspVEvKd/SOfKIN5gnWm6JoQdkvn
NoyGXZD4eBgYebULjyySkFcUxkM0Yr9Dv8jDM9rZ12Yk+3im29k5nFTF4d3XOv3A
tDbR5CFD2f8nBSMi1+Y2aJOd7UXlBN7RLYRbPlIBmGXQFvwh9vee1gY0J6U/54sI
iaFv1QKBgQDxIWYt2YJ5woQWkdGEFlQS5iFvq0Z4fNllTzCr0xNcaxeb13ksPtxj
5Q2DGYFjE2+UI/uJwFVht56yrJzJdCaFdVL2ImpwmCCbzPj+Nk24WmLfPzXTvIB4
26zYrT5TFk3D+iv9pNzFBq74vZo9GCVkGWv3Zk6SY3KiLnft/uT4gA==
-----END RSA PRIVATE KEY-----`)
)
