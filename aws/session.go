package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/blend/go-sdk/ex"
)

// MustNewSession creates a new aws session from a config and panics on error.
func MustNewSession(cfg Config) *session.Session {
	session, err := NewSession(cfg)
	if err != nil {
		panic(err)
	}
	return session
}

// NewSession creates a new aws session from a config.
func NewSession(cfg Config) (*session.Session, error) {
	if cfg.IsZero() {
		session, err := session.NewSession()
		if err != nil {
			return nil, ex.New(err)
		}
		return session, nil
	}

	awsConfig := &aws.Config{
		Region:      aws.String(cfg.Region),
		Credentials: credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, cfg.SecurityToken),
	}
	session, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, ex.New(err)
	}
	return session, nil
}
