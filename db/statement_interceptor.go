package db

// StatementInterceptor is an interceptor for statements.
type StatementInterceptor func(statementID, statement string) string
