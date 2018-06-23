package logger

const (
	// FlagAll is a special flag that allows all events to fire.
	FlagAll Flag = "all"
	// FlagNone is a special flag that allows no events to fire.
	FlagNone Flag = "none"

	// Fatal fires for fatal errors and is an alias to `Fatal`.
	Fatal Flag = "fatal"
	// Error fires for errors that are severe enough to log but not so severe as to abort a process.
	Error Flag = "error"
	// Warning fires for warnings.
	Warning Flag = "warning"
	// Debug fires for debug messages.
	Debug Flag = "debug"
	// Info fires for informational messages (app startup etc.)
	Info Flag = "info"
	// Silly is for when you just need to log something weird.
	Silly Flag = "silly"

	// WebRequestStart is an event flag.
	WebRequestStart Flag = "web.request.start"
	// WebRequest is an event flag.
	WebRequest Flag = "web.request"

	// Audit is an event flag.
	Audit Flag = "audit"

	// Query is a logging flag.
	Query Flag = "db.query"

	// FrontendAction flag is a frontend action
	FrontendAction Flag = "frontend.action"
)

// Flag represents an event type that can be enabled or disabled.
type Flag string

// AsFlags casts a variadic list of strings to an array of Flag.
func AsFlags(values ...string) (output []Flag) {
	for _, v := range values {
		output = append(output, Flag(v))
	}
	return
}

// AsStrings casts a variadic list of flags to an array of strings.
func AsStrings(values ...Flag) (output []string) {
	for _, v := range values {
		output = append(output, string(v))
	}
	return
}
