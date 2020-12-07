package oauth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"time"

	"github.com/blend/go-sdk/jwt"
	"github.com/blend/go-sdk/uuid"
)

const (
	pk0pem = "-----BEGIN RSA PRIVATE KEY-----\nMIIEoQIBAAKCAQEAy4zZIH5mtomtQisfKhhr79LMbJqtWJWRDJytoxWb3h0gNniz\n2uxLVBqIyIibiqDr9dG0waQXlZgWip7WnL0Rip+UjGN+i0jhHAqNfjpK0KVI/epa\nFeL3rP+jfTZXOedUwt7kAuxCI07Dokqyarm7WpAEaShkd8ZuQ1KsQ9ZblKrw77uO\n6fB34npf+2Lahi++P0FFpgHNW2vfSc0PrXoGy27DkWDMnHGiElm8VX9nwOGI+JN1\n3RdQzTSF4VNoXjRynLNdPt4XKLQS0HFW3kQHogi0uf0KG6sZpLBU5KgORFD7ScXm\nAd402NH7qYdzggdhaSMsTZ16RpcdBtaX/KF+GwIDAQABAoIBABrbR56w7s5w1epg\nFCmStVMcRhqiQfLpMQ0v8v0Mkdc5kpF9VYWyHbJIGfoThCpDVz7E34uZIf975KV/\nlaNyksjui0QGsKCiCgmQHuEjwdFLrZjK/f3bR4CM7j5MGDAspJNdo0n7cDKGZuuX\n3XiVbvHhBKP3T2I6TTwWwWHl+4le1FNRvcPPK1V4SXrdknuZl1Rbzz2xt2AC4hZR\niiWhA893pXnai4IbjHwcOUGRKT6i63TNjsik/o0ANJruoyVfdMyRBnTQ3cUUM1/k\n+qWEki1fBHLVcSrhCdAuq0rezCZnTjd6Z2t97XOYMm/t9ak0f8ZvtLntomLxv8Mv\n5BJel2ECgYEA5wFwGILXWWmPmHPltTRd6/v8ZDjFpLWJ/xfdnp7XfDwXks5HbYV8\nsC6A2fPotYuHFM4m3kP3HteKqRUYndrG7Aj1CYFoa3GKE1ycepj6+SPahabgRrxZ\nrlQQH2rTwQYxl7WGJf2cHGa71YYJZNmXKtDY1eR49l0NPkELV5r93usCgYEA4ZLu\nBkyph/fJ3AMIGrc3fTyiHt+DmFuu1S0AkcvH1NOCRLOVyxJdwswr26q7Wn4tqF5h\ntaRSwXzoyUEESOZK+IdQdbtoRYxLw7SZITyxfBg+Ds8U8n65Occ2NOW9PQNIX9rt\nzf5bO3u4AIklfUKAWVL/ufw+Dj4Hb7lH8Vsx8ZECgYBWEogU6fOZgiaZ9F0btmZk\nfmCdazXhWC8R2G+gIalCxhU2gxvEKB+8eadTDnmf41wymVmMKaDTYhZtR8oDTzgd\nTH0YzJn+prB+5Fv9pjClUgGjGPmqAZYcyX+0ZRZ/bnJeB6nzT9qyDmlgdu/bHuQf\ndO/GSrnzedpsXsn+G2cKfwJ/XkgLNJbWRP5MYKjjukbZ6n5tRHonhobLjE5C7q09\n2LaOvChTc405ozGzIx05MZmLe9P3AvSrojOTGIsUP2QB8d6cwpiR/H+nKyVQ25OC\nm1uGlKn5F2HgCUY6YeGkNtwoY+gdfPvTJgmP3ql0AebJvovyVsoXJdzHPusyJq73\nMQKBgQDS7xbSGoaDNwhIwnncaYuCDdi8AcQ1EtwCg4YdQ0HdpuukLBMBY0w00Ueq\npirBbNHduauDpAnStTHYEbTnS1p/VV7UkKjoTDpw06h8l56UToU7wmhdZ17LuBDx\n+FH/GEF2qMkwGXrViVnaEL9ZmhC53s6qepmKlspGLVcUt4MBdA==\n-----END RSA PRIVATE KEY-----"
	pk1pem = "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAz3/uFaNQPyWNVms4+9G7q/pWoTpGfEeSD8FTuA19bGGv1TaJ\nPV55tKhoKLNsKp1C5E/F+OmXgs30S3VoBDVxUuVTROrD0WygwDNO/O4OzrP3mbWT\nzCO6BH8x3qholqobaRgHML1WIhA+wAIUfXNwhuIkVtusvcaPd3pM3IMTrC5tC/iS\nFlbuZkN/+T2seMYHCQ4OIvzee2+kJYBcBHs1nwKg8NStN9f3YjGkauqwqA0dV8JG\nEcHe0qercwOtr/k1O47Mo2+TyV7fzc8JEHsfiGxaCyiZw1ocjFZGQqTXD74Ufvqx\naWcph/NChqxEt6kpNdJEyB92nPSLud7WT3N1YwIDAQABAoIBAQCxMlCBDewTYOAn\n8nzBH0QjAy9Dk95pdz0WU0RJIsv+6BUeAOqGC83nJwF78Gzon09mZXFstR57x6Fd\nZy+imHjkD45ihhEfIKLOP4KuoCTpA+rnypYieEf8WxqdSDe4oh+ySaCqUKXjhPfx\nRFV3JEPuC+R4gDQuBAi0QS6uCQmduKX/KpG7SGNCmLE1Qc1kNHJUMbmAj9YUp1tN\nGndCoCEgCo8Xvnjs/rFi9rpUFnBfrSTosxm+4YN9G9SHX1UG2yPiEviU/uJeFXC3\npjUJe/AWHA746TFEJYitQtN693tx5ZpOq0gX7Oyn/sse9uqgKvL/r67tmcHxq1Jn\n6b1rsR5xAoGBAOz5nto+xMN52ewbSNH53HW2XO1c0y1eUYrGHLSJFSeYPjDLKHrV\nFqPCpYb0SbNE9GiuEwWqdG08KhTc17IC1BwieHmCkw53etLgQKsmQvloiONgESDI\nUogPr8ENPZkRyrLh7WWKITpayASAkh6k4/pQ1JLrZlrh9h5dbNuENr/dAoGBAOAo\nhgwXKTlIMIK9zefe5fOcUW2Jfl3aD/Jcl0bsuQqKRJ7UgoqBdYC8UYzZyPUECZ4H\nxqRb8Y4sfbV0/bxuuXtdW571ru/pqfH9c6sstvT/VCO7yWriyrPmH/GeRGKD6k4d\n8h9CTNIEKLsCsuIh26DHuNv8GVexFuhyhQnyPlY/AoGADt9HoejIjoAKNjAsLMli\nlZyhTmBB/JnrwirWyFnGExsR5BwL6VGQPyzLGKIiMfcE48Dw/q0I64YYGgEWJFzb\nFPzw1KdmNUU4Vx2t0U/wahiuZp6z1Hvd+h4J6LK9B+s+7mURcgruNOxXmzi6cuPk\nTuRdwu61GMUPni4807YDfZ0CgYB53q5afnEkOpJdUrJTAUXGN9OxmRJCFl+bJin9\nHpDQITKDpAhBI3duAXTY/kMaqxJLf/DIxVBEOv4xnKSjQRPI9Y3tk6eDumdyMJkl\nlI74DqWBNASi/yCzxEbTx3dolE3cIL3VruczO52lZyc4eK3+8PcZayugGKDaygB4\n0uJ/YwKBgHRUFlDpNLDQ0UPoX0doDNuKIb/yeV0yz7aSSfr5cAkEFlHr+b+25lON\nMez+mWQgoKue3qQYK5z+aP/1ebHN3OdTSRjSh3sqgqTM0hhcd2K66Wfd3v9RyysG\ne0I+hQENst5tk7HgkqN4IZ3ha1iw/1oyMHFxwY6lE40qBcmAPi2f\n-----END RSA PRIVATE KEY-----"
	pk2pem = "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEAw2dRbUsAtrbqaTWN9Ekt2Zq6fZfRU+YkOnxIDKt4USsdI8Gq\nm58lB+wJ4RXzR2YPTa+YFXALoDvE/iKu+BpJQSj3WPoV2zvLM7InXcUKUK2+qDmu\na2etu8otFjA92WwcIxMBuPI34rr99UzZiUO2Yf7E4zYSWunGmCuiUOkTjrldjqpk\nCDGt9BqHC8JdnKCBn7KnMX4DisVdOC1d4NtpQTarfUl1yCf1da737GEjzGjIjQkc\nSCXTaGOIEszEeBYp6Elu/H1Ay8u7RkRvU7ufRRn9sYrr+ek1R3CwKw/mT2Ot9cYe\nFZGabNS84UgeqeMtvlSmTv2/kWX7ITEn/p6lIQIDAQABAoIBAHIMXIdIznrWWgzc\nGCVrjNpEJ/Lj6GZqndyQ61CRyCC/5DsZbyVzhp6QEtgQArU6iVYTVdW1VuPH3tth\njPP8C6N/cJa7KIST6q8anUVqmvGp5uyy9e10Tv+bKiOYNpEvO2DxWAEFRr8L2uwQ\nVat7HPknRO1EgwQTDDmGxi8pSqPy4bsyLA8lW7uDlrVQdqCGau6uL9Tc/Jq71uP7\n5z+u/tOJZY5SQ2E16RKp89HpQYrOkFvckdpSJd78XdX0jhqiIeH+DQQIssJvCK8q\nslIcCVkoSuuvYzL465EZCqGGdzQv66ywEIVd3HaUkBZdybk4dghLGu5vuyrdJ/Bm\n7RHgYi0CgYEA580hQIxiHa5mloS9RSjdqWI8K6qp+KvDSVinmr255NsExXM5y3fX\n+NwsaAcaWFVlYGdYlJVdTR5I0X9BYUEuTLM+hf5M8QOG95KR3vXnHpfl745NfbDC\nVcDeOIc/XzWUAfCku+US2ijGfP1bgMX1UQnYHz1BUt9BC14N12Oi8IcCgYEA1812\nblT6MkuOiiTl5RRgm2xnd69MxI2wizQ503hTmEJmtaWQ14M6Yz9YvkbH8jwjMhuo\nZMwI3A4AvYH8sYfKlJI5gRNoAldHOf/nbQEB7d9QQel6/0QSQ05ZlO4FmG/kHLhO\nsPjll1nOeki3WXb8NZ5K8LxVZOUc/K6Xq3xr7xcCgYEAuRBhwuoRj4bkqrlRbvzg\nc9JVHbvEth9T66QXNAjTeG6QEaAb/WEyEaKe5XL+SpXrORtpcj8J3X8XPgMuTJpA\nf8X/XfUYsrdRMylWwr5qhldZoXdoULglf1dbU6BPLRFWmHHq44RRF9HEHpgcTOQ/\nJjMI1HAQTjyl7pBp1pPay9MCgYEAiNE8uqqpjWWV00OddWU78o4B80FyvFLQkRDl\ncIsjBK9kitmTQO9z/yRUUR5y+cLi1YvvcShinZFLKtrUqIFdEGC8kHcLRCCtiboS\nsWsoG/Wu3nr2fgxcP8vWw7M8XO7jgsnfKhhDB3fqjmC3zcLAGAZpoMLmqPcRL6pJ\ngnF5xLUCgYAI1rar19R1uNXXiH4JVtUk57c643sMRN/09J+nbyT7xkWtPtUujNeY\nxj3iakptk3j/mW7H5qDE9p3kvDAI7Xp3lg/0t6pOWIajuRZueI/36pMSUpkUqO8F\nBua3vojOESMkzFaI5oXLjzTZ1Om+BiqeOYgorA7FKjcYwFRyilcspg==\n-----END RSA PRIVATE KEY-----"
)

func createJWK(pk *rsa.PrivateKey) jwt.JWK {
	eBytes := big.Int{}
	eBytes.SetInt64(int64(pk.PublicKey.E))
	e := base64.RawURLEncoding.EncodeToString(eBytes.Bytes())
	n := base64.RawURLEncoding.EncodeToString(pk.PublicKey.N.Bytes())
	return jwt.JWK{
		KID: uuid.V4().String(),
		ALG: "RS256",
		KTY: "RSA",
		USE: "sig",
		E:   e,
		N:   n,
	}
}

func createCodeResponse(aud, keyID string, pk *rsa.PrivateKey) ([]byte, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, &GoogleClaims{
		StandardClaims: jwt.StandardClaims{
			Audience:  aud,
			ExpiresAt: time.Now().UTC().AddDate(0, 0, 1).Unix(),
			IssuedAt:  time.Now().UTC().Unix(),
			Issuer:    GoogleIssuer,
		},
		HD:            "test.blend.com",
		Email:         "example-string@test.blend.com",
		EmailVerified: "true",
	})
	jwtToken.Header["kid"] = keyID
	jwtTokenSigned, err := jwtToken.SignedString(pk)
	if err != nil {
		return nil, err
	}
	type tokenJSON struct {
		AccessToken  string `json:"access_token"`
		IDToken      string `json:"id_token"`
		TokenType    string `json:"token_type"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		Expires      int    `json:"expires"`
	}
	return json.Marshal(tokenJSON{
		AccessToken: "test_access_token",
		IDToken:     jwtTokenSigned,
	})
}
