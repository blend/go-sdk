package tracing

// These constants are mostly lifted from the datadog/tracing/ext tag values.
const (
	// TagKeyEnvironment is the environment (web, dev, etc.)
	TagKeyEnvironment = "env"
	// TagKeySpanType defines the Span type (web, db, cache).
	TagKeySpanType = "span.type"
	// TagKeyServiceName defines the Service name for this Span.
	TagKeyServiceName = "service.name"
	// TagKeyResourceName defines the Resource name for the Span.
	TagKeyResourceName = "resource.name"
	// TagKeyPID is the pid of the traced process.
	TagKeyPID = "system.pid"
	// TagKeyError is the error tag key. It is usually of type `error`.
	TagKeyError = "error"
	// TagKeyErrorMessage is the error message tag key.
	TagKeyErrorMessage = "error.message"
	// TagKeyErrorStack is the error stack tag key.
	TagKeyErrorStack = "error.stack"
	// TagKeyHTTPMethod is the verb on the request.
	TagKeyHTTPMethod = "http.method"
	// TagKeyHTTPCode is the result status code.
	TagKeyHTTPCode = "http.status_code"
	// TagKeyHTTPURL is the url of the request (typically the raw path).
	TagKeyHTTPURL = "http.url"
	// TagKeyDBApplication is the application that uses a database.
	TagKeyDBApplication = "db.application"
	// TagKeyDBName is the database name.
	TagKeyDBName = "db.name"
	// TagKeyDBRowsAffected is the number of rows affected.
	TagKeyDBRowsAffected = "db.rows_affected"
	// TagKeyDBUser is the user on the database connection.
	TagKeyDBUser = "db.user"
	// TagKeyJobName is the job name.
	TagKeyJobName = "job.name"
	// TagKeyGRPCRole is the grpc role (i.e. client or server).
	TagKeyGRPCRole = "grpc.role"
	// TagKeyGRPCMethod is the grpc method.
	TagKeyGRPCMethod = "grpc.method"
	// TagKeyGRPCCode is the grpc result code.
	TagKeyGRPCCode = "grpc.code"
	// TagKeyGRPCAuthority is the grpc authority.
	TagKeyGRPCAuthority = "grpc.authority"
	// TagKeyGRPCUserAgent is the grpc user-agent.
	TagKeyGRPCUserAgent = "grpc.user_agent"
	// TagKeyGRPCContentType is the grpc content type.
	TagKeyGRPCContentType = "grpc.content_type"
	// TagSecretsOperation is the operation being performed in the secrets API
	TagSecretsOperation = "secrets.operation"
	// TagSecretsMethod is the http method being hit on the vault API
	TagSecretKey = "secrets.key"
	// TagKeyOAuthUsername defines the oauth Username name for the Span.
	TagKeyOAuthUsername = "oauth.username"
)

// Operations are actions represented by spans.
const (
	OperationHTTPRouteLookup = "http.route_lookup"
	// OperationHTTPRequest is the http request tracing operation name.
	OperationHTTPRequest = "http.request"
	// OperationHTTPRender is the operation name for rendering a server side view.
	OperationHTTPRender = "http.render"
	// OperationDBPing is the db ping tracing operation.
	OperationSQLPing = "sql.ping"
	// OperationDBPrepare is the db prepare tracing operation.
	OperationSQLPrepare = "sql.prepare"
	// OperationDBQuery is the db query tracing operation.
	OperationSQLQuery = "sql.query"
	// OperationJob is a job operation.
	OperationJob = "job"
	// OperationRPC is an rpc operation.
	OperationRPC = "rpc"
	// OperationVaultAPI is a call to the vault API
	OperationVaultAPI = "vault.api.request"
)

// Span types have similar behaviour to "app types" and help categorize
// traces in the Datadog application. They can also help fine grain agent
// level bahviours such as obfuscation and quantization, when these are
// enabled in the agent's configuration.
const (
	// SpanTypeWeb marks a span as an HTTP server request.
	SpanTypeWeb = "web"
	// SpanTypeHTTP marks a span as an HTTP client request.
	SpanTypeHTTP = "http"
	// SpanTypeSQL marks a span as an SQL operation. These spans may
	// have an "sql.command" tag.
	SpanTypeSQL = "sql"
	// SpanTypeCassandra marks a span as a Cassandra operation. These
	// spans may have an "sql.command" tag.
	SpanTypeCassandra = "cassandra"
	// SpanTypeRedis marks a span as a Redis operation. These spans may
	// also have a "redis.raw_command" tag.
	SpanTypeRedis = "redis"
	// SpanTypeMemcached marks a span as a memcached operation.
	SpanTypeMemcached = "memcached"
	// SpanTypeMongoDB marks a span as a MongoDB operation.
	SpanTypeMongoDB = "mongodb"
	// SpanTypeElasticSearch marks a span as an ElasticSearch operation.
	// These spans may also have an "elasticsearch.body" tag.
	SpanTypeElasticSearch = "elasticsearch"
	// SpanTypeJob is a span type used by cron jobs.
	SpanTypeJob = "job"
	// SpanTypeGRPC is a span type used by grpc services.
	SpanTypeGRPC = "grpc"
	// SpanTypeVault is a span type used by go-sdk/secrets calls to vault
	SpanTypeVault = "vault"
)

// Priority is a hint given to the backend so that it knows which traces to reject or kept.
// In a distributed context, it should be set before any context propagation (fork, RPC calls) to be effective.
const (
	// PriorityUserReject informs the backend that a trace should be rejected and not stored.
	// This should be used by user code overriding default priority.
	PriorityUserReject = -1

	// PriorityAutoReject informs the backend that a trace should be rejected and not stored.
	// This is used by the builtin sampler.
	PriorityAutoReject = 0

	// PriorityAutoKeep informs the backend that a trace should be kept and not stored.
	// This is used by the builtin sampler.
	PriorityAutoKeep = 1

	// PriorityUserKeep informs the backend that a trace should be kept and not stored.
	// This should be used by user code overriding default priority.
	PriorityUserKeep = 2
)
