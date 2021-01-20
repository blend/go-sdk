package selector

const (
	// At is a common rune.
	At = rune('@')
	// Colon is a common rune.
	Colon = rune(':')
	// Dash is a common rune.
	Dash = rune('-')
	// Underscore  is a common rune.
	Underscore = rune('_')
	// Dot is a common rune.
	Dot = rune('.')
	// ForwardSlash is a common rune.
	ForwardSlash = rune('/')
	// BackSlash is a common rune.
	BackSlash = rune('\\')
	// BackTick is a common rune.
	BackTick = rune('`')
	// Bang is a common rune.
	Bang = rune('!')
	// Comma is a common rune.
	Comma = rune(',')
	// OpenBracket is a common rune.
	OpenBracket = rune('[')
	// OpenParens is a common rune.
	OpenParens = rune('(')
	// OpenCurly is a common rune.
	OpenCurly = rune('{')
	// CloseBracket is a common rune.
	CloseBracket = rune(']')
	// CloseParens is a common rune.
	CloseParens = rune(')')
	// Equal is a common rune.
	Equal = rune('=')
	// Space is a common rune.
	Space = rune(' ')
	// Tab is a common rune.
	Tab = rune('\t')
	// Tilde is a common rune.
	Tilde = rune('~')
	// CarriageReturn is a common rune.
	CarriageReturn = rune('\r')
	// NewLine is a common rune.
	NewLine = rune('\n')
)

const (
	// OpEquals is an operator.
	OpEquals = "="
	// OpDoubleEquals is an operator.
	OpDoubleEquals = "=="
	// OpNotEquals is an operator.
	OpNotEquals = "!="
	// OpIn is an operator.
	OpIn = "in"
	// OpNotIn is an operator.
	OpNotIn = "notin"
)

const (
	// ErrInvalidOperator is returned if the operator is invalid.
	ErrInvalidOperator Error = "invalid operator"
	// ErrInvalidSelector is returned if there is a structural issue with the selector.
	ErrInvalidSelector Error = "invalid selector"
	// ErrLabelKeyEmpty indicates a key is empty.
	ErrLabelKeyEmpty Error = "label key empty"
	// ErrLabelKeyTooLong indicates a key is too long.
	ErrLabelKeyTooLong Error = "label key too long"
	// ErrLabelKeyDNSSubdomainEmpty indicates a key's "dns" subdomain is empty, i.e. it is in the form `/foo`
	ErrLabelKeyDNSSubdomainEmpty Error = "label key dns subdomain empty"
	// ErrLabelKeyDNSSubdomainTooLong indicates a key's "dns" subdomain is too long.
	ErrLabelKeyDNSSubdomainTooLong Error = "label key dns subdomain too long; must be less than 253 characters"
	// ErrLabelValueTooLong indicates a value is too long.
	ErrLabelValueTooLong Error = "label value too long; must be less than 63 characters"
	// ErrLabelInvalidCharacter indicates a value contains characters
	ErrLabelInvalidCharacter Error = `label contains invalid characters, regex used: ([A-Za-z0-9_-.])`
	// ErrLabelKeyInvalidDNSSubdomain indicates a key contains characters
	ErrLabelKeyInvalidDNSSubdomain Error = `label key dns subdomain contains invalid dns characters, regex used: ([a-z0-9-.])`

	// MaxLabelKeyDNSSubdomainLen is the maximum dns prefix length.
	MaxLabelKeyDNSSubdomainLen = 253
	// MaxLabelKeyLen is the maximum key length.
	MaxLabelKeyLen = 63
	// MaxLabelValueLen is the maximum value length.
	MaxLabelValueLen = 63
)

var (
	// MaxLabelKeyTotalLen is the maximum total key length.
	MaxLabelKeyTotalLen = MaxLabelKeyDNSSubdomainLen + MaxLabelKeyLen + 1
)
