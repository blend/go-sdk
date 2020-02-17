package email

import "github.com/blend/go-sdk/configutil"

var (
	_ configutil.Resolver = (*SMTPSender)(nil)
	_ configutil.Resolver = (*SMTPPlainAuth)(nil)
)
