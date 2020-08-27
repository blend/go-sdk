# Capping Query Running Time with Statement Timeout

## Prerequisites

Set some common environment variables (we `export` here to make running
the Go script a bit simpler, but these can be local in a shell or local
to a command)

```
export DB_HOST=localhost
export DB_PORT=30071
export DB_USER=superuser
export DB_NAME=superuser_db
export DB_PASSWORD=testpassword_superuser
export DB_SSLMODE=disable
```

and make sure a local `postgres` server is running

```
docker run \
  --detach \
  --hostname "${DB_HOST}" \
  --publish "${DB_PORT}:5432" \
  --name dev-postgres-statement-timeout \
  --env "POSTGRES_DB=${DB_NAME}" \
  --env "POSTGRES_USER=${DB_USER}" \
  --env "POSTGRES_PASSWORD=${DB_PASSWORD}" \
  postgres:10.6-alpine
```

## Let `postgres` Cancel Via `statement_timeout`

```
$ go run .
0.000053 ==================================================
0.000087 Configured statement timeout: 10ms
0.000091 Configured pg_sleep:          200ms
0.000110 Configured context timeout:   400ms
0.000141 ==================================================
0.014139 DSN="postgres://superuser:testpassword_superuser@localhost:30071/superuser_db?sslmode=disable&statement_timeout=10ms"
0.014147 ==================================================
0.015832 statement_timeout=10ms
0.038094 ==================================================
0.038103 Error(s):
0.038113 - Message: "SELECT id, pg_sleep(0.200000) FROM might_sleep WHERE id = 1337;"
0.038139 - &pq.Error{Severity:"ERROR", Code:"57014", Message:"canceling statement due to statement timeout", Detail:"", Hint:"", Position:"", InternalPosition:"", InternalQuery:"", Where:"", Schema:"", Table:"", Column:"", DataTypeName:"", Constraint:"", File:"postgres.c", Line:"2996", Routine:"ProcessInterrupts"}
```

From [Appendix A. PostgreSQL Error Codes][1]:

```
Class 57 - Operator Intervention
---------+-----------------------
   57014 | query_canceled
```

## Cancel Query via Go `context` Cancelation

```
$ VIA_GO_CONTEXT=true go run .
0.000072 ==================================================
0.000105 Configured statement timeout: 10s
0.000117 Configured pg_sleep:          200ms
0.000120 Configured context timeout:   100ms
0.000146 ==================================================
0.013925 DSN="postgres://superuser:testpassword_superuser@localhost:30071/superuser_db?sslmode=disable&statement_timeout=10000ms"
0.013934 ==================================================
0.015556 statement_timeout=10s
0.108807 ==================================================
0.108817 Error(s):
0.108831 - Message: "SELECT id, pg_sleep(0.200000) FROM might_sleep WHERE id = 1337;"
0.108863 - &pq.Error{Severity:"ERROR", Code:"57014", Message:"canceling statement due to user request", Detail:"", Hint:"", Position:"", InternalPosition:"", InternalQuery:"", Where:"", Schema:"", Table:"", Column:"", DataTypeName:"", Constraint:"", File:"postgres.c", Line:"3026", Routine:"ProcessInterrupts"}
```

## `psql` Does **NOT** Support `statement_timeout` in DSN

```
$ psql "postgres://superuser:testpassword_superuser@localhost:30071/superuser_db?sslmode=disable&statement_timeout=10ms"
psql: error: could not connect to server: invalid URI query parameter: "statement_timeout"
$ psql "postgres://superuser:testpassword_superuser@localhost:30071/superuser_db?sslmode=disable"
psql (12.4, server 10.6)
Type "help" for help.

superuser_db=# \q
```

## Clean Up

```
docker rm --force dev-postgres-statement-timeout
```

[1]: https://www.postgresql.org/docs/10/errcodes-appendix.html
