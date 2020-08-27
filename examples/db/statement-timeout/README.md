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
0.000055 ==================================================
0.000090 Configured statement timeout: 10ms
0.000095 Configured pg_sleep:          200ms
0.000098 Configured context timeout:   400ms
0.000124 ==================================================
0.015619 DSN="postgres://superuser:testpassword_superuser@localhost:30071/superuser_db?sslmode=disable&statement_timeout=10ms"
0.015629 ==================================================
0.016960 statement_timeout=10ms
0.024797 ==================================================
0.024812 Starting query
0.036470 ***
0.036478 Error(s):
0.036488 - Message: "SELECT id, pg_sleep(0.200000) FROM might_sleep WHERE id = 1337;"
0.036529 - &pq.Error{Severity:"ERROR", Code:"57014", Message:"canceling statement due to statement timeout", Detail:"", Hint:"", Position:"", InternalPosition:"", InternalQuery:"", Where:"", Schema:"", Table:"", Column:"", DataTypeName:"", Constraint:"", File:"postgres.c", Line:"2996", Routine:"ProcessInterrupts"}
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
0.000081 ==================================================
0.000116 Configured statement timeout: 10s
0.000120 Configured pg_sleep:          200ms
0.000133 Configured context timeout:   100ms
0.000163 ==================================================
0.014563 DSN="postgres://superuser:testpassword_superuser@localhost:30071/superuser_db?sslmode=disable&statement_timeout=10000ms"
0.014575 ==================================================
0.016120 statement_timeout=10s
0.023707 ==================================================
0.023729 Starting query
0.106258 ***
0.106272 Error(s):
0.106309 - Message: "SELECT id, pg_sleep(0.200000) FROM might_sleep WHERE id = 1337;"
0.106341 - &pq.Error{Severity:"ERROR", Code:"57014", Message:"canceling statement due to user request", Detail:"", Hint:"", Position:"", InternalPosition:"", InternalQuery:"", Where:"", Schema:"", Table:"", Column:"", DataTypeName:"", Constraint:"", File:"postgres.c", Line:"3026", Routine:"ProcessInterrupts"}
```

## `psql` Does **NOT** Support `statement_timeout` in DSN

See `libpq` [Parameter Key Words][2]

```
$ psql "postgres://superuser:testpassword_superuser@localhost:30071/superuser_db?sslmode=disable&statement_timeout=10ms"
psql: error: could not connect to server: invalid URI query parameter: "statement_timeout"
$ psql "postgres://superuser:testpassword_superuser@localhost:30071/superuser_db?sslmode=disable"
psql (12.4, server 10.6)
Type "help" for help.

superuser_db=# \q
$
$
$ psql "user=superuser password=testpassword_superuser host=localhost port=30071 dbname=superuser_db sslmode=disable statement_timeout=10ms"
psql: error: could not connect to server: invalid connection option "statement_timeout"
$ psql "user=superuser password=testpassword_superuser host=localhost port=30071 dbname=superuser_db sslmode=disable"
psql (12.4, server 10.6)
Type "help" for help.

superuser_db=# \q
```

## Clean Up

```
docker rm --force dev-postgres-statement-timeout
```

[1]: https://www.postgresql.org/docs/10/errcodes-appendix.html
[2]: https://www.postgresql.org/docs/10/libpq-connect.html#LIBPQ-PARAMKEYWORDS
