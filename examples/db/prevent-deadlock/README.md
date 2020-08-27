# Preventing Deadlock with Lock Timeout

## Prerequisites

Set some common environment variables (we `export` here to make running
the Go script a bit simpler, but these can be local in a shell or local
to a command)

```
export DB_HOST=localhost
export DB_PORT=28007
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
  --name dev-postgres-prevent-deadlock \
  --env "POSTGRES_DB=${DB_NAME}" \
  --env "POSTGRES_USER=${DB_USER}" \
  --env "POSTGRES_PASSWORD=${DB_PASSWORD}" \
  postgres:10.6-alpine
```

## Let `postgres` Cancel Via `lock_timeout`

```
$ go run .
0.000055 ==================================================
0.000085 Configured lock timeout:      10ms
0.000089 Configured context timeout:   600ms
0.000091 Configured transaction sleep: 200ms
0.000114 ==================================================
0.014372 DSN="postgres://superuser:testpassword_superuser@localhost:28007/superuser_db?lock_timeout=10ms&sslmode=disable"
0.014381 ==================================================
0.015569 lock_timeout=10ms
0.026958 ==================================================
0.026981 Starting transactions
0.036223 Transactions opened
0.261793 ***
0.261803 Error(s):
0.261852 - &pq.Error{Severity:"ERROR", Code:"55P03", Message:"canceling statement due to lock timeout", Detail:"", Hint:"", Position:"", InternalPosition:"", InternalQuery:"", Where:"while updating tuple (0,2) in relation \"might_deadlock\"", Schema:"", Table:"", Column:"", DataTypeName:"", Constraint:"", File:"postgres.c", Line:"2989", Routine:"ProcessInterrupts"}
0.261862 - &pq.Error{Severity:"ERROR", Code:"55P03", Message:"canceling statement due to lock timeout", Detail:"", Hint:"", Position:"", InternalPosition:"", InternalQuery:"", Where:"while updating tuple (0,1) in relation \"might_deadlock\"", Schema:"", Table:"", Column:"", DataTypeName:"", Constraint:"", File:"postgres.c", Line:"2989", Routine:"ProcessInterrupts"}
0.261870 - &errors.errorString{s:"pq: Could not complete operation in a failed transaction"}
0.261874 - &errors.errorString{s:"pq: Could not complete operation in a failed transaction"}
```

From [Appendix A. PostgreSQL Error Codes][1]:

```
Class 55 - Object Not In Prerequisite State
---------+----------------------------------
   55P03 | lock_not_available
```

## Force a Deadlock

By allowing the Go context to stay active for **very** long, we can allow
Postgres to detect

```
$ FORCE_DEADLOCK=true go run .
0.000044 ==================================================
0.000068 Configured lock timeout:      10s
0.000071 Configured context timeout:   10s
0.000073 Configured transaction sleep: 200ms
0.000089 ==================================================
0.011839 DSN="postgres://superuser:testpassword_superuser@localhost:28007/superuser_db?lock_timeout=10000ms&sslmode=disable"
0.011850 ==================================================
0.013332 lock_timeout=10s
0.022643 ==================================================
0.022659 Starting transactions
0.030515 Transactions opened
1.245005 ***
1.245016 Error(s):
1.245053 - &pq.Error{Severity:"ERROR", Code:"40P01", Message:"deadlock detected", Detail:"Process 347 waits for ShareLock on transaction 845; blocked by process 346.\nProcess 346 waits for ShareLock on transaction 846; blocked by process 347.", Hint:"See server log for query details.", Position:"", InternalPosition:"", InternalQuery:"", Where:"while updating tuple (0,1) in relation \"might_deadlock\"", Schema:"", Table:"", Column:"", DataTypeName:"", Constraint:"", File:"deadlock.c", Line:"1140", Routine:"DeadLockReport"}
1.245063 - &errors.errorString{s:"pq: Could not complete operation in a failed transaction"}
```

From [Appendix A. PostgreSQL Error Codes][1]:

```
Class 40 - Transaction Rollback
---------+----------------------
   40P01 | deadlock_detected
```

## Go `context` Cancelation In Between Queries in a Transaction

```
$ BETWEEN_QUERIES=true go run .
0.000051 ==================================================
0.000082 Configured lock timeout:      10s
0.000086 Configured context timeout:   100ms
0.000089 Configured transaction sleep: 200ms
0.000110 ==================================================
0.013163 DSN="postgres://superuser:testpassword_superuser@localhost:28007/superuser_db?lock_timeout=10000ms&sslmode=disable"
0.013176 ==================================================
0.014402 lock_timeout=10s
0.025375 ==================================================
0.025401 Starting transactions
0.032665 Transactions opened
0.236497 ***
0.236519 Error(s):
0.236575 - context.deadlineExceededError{}
0.236587 - Context cancel in between queries
0.236591 - context.deadlineExceededError{}
0.236615 - Context cancel in between queries
0.236636 - &errors.errorString{s:"sql: transaction has already been committed or rolled back"}
0.236643 - &errors.errorString{s:"sql: transaction has already been committed or rolled back"}
```

## Cancel "Stuck" Deadlock via Go `context` Cancelation

```
$ DISABLE_LOCK_TIMEOUT=true go run .
0.000053 ==================================================
0.000084 Configured lock timeout:      10s
0.000088 Configured context timeout:   600ms
0.000091 Configured transaction sleep: 200ms
0.000113 ==================================================
0.014431 DSN="postgres://superuser:testpassword_superuser@localhost:28007/superuser_db?lock_timeout=10000ms&sslmode=disable"
0.014442 ==================================================
0.016239 lock_timeout=10s
0.026890 ==================================================
0.026903 Starting transactions
0.036405 Transactions opened
0.612462 ***
0.612474 Error(s):
0.612539 - &pq.Error{Severity:"ERROR", Code:"57014", Message:"canceling statement due to user request", Detail:"", Hint:"", Position:"", InternalPosition:"", InternalQuery:"", Where:"while updating tuple (0,2) in relation \"might_deadlock\"", Schema:"", Table:"", Column:"", DataTypeName:"", Constraint:"", File:"postgres.c", Line:"3026", Routine:"ProcessInterrupts"}
0.612556 - &pq.Error{Severity:"ERROR", Code:"57014", Message:"canceling statement due to user request", Detail:"", Hint:"", Position:"", InternalPosition:"", InternalQuery:"", Where:"while updating tuple (0,1) in relation \"might_deadlock\"", Schema:"", Table:"", Column:"", DataTypeName:"", Constraint:"", File:"postgres.c", Line:"3026", Routine:"ProcessInterrupts"}
0.612569 - &errors.errorString{s:"sql: transaction has already been committed or rolled back"}
0.612578 - &errors.errorString{s:"sql: transaction has already been committed or rolled back"}
```

From [Appendix A. PostgreSQL Error Codes][1]:

```
Class 57 - Operator Intervention
---------+-----------------------
   57014 | query_canceled
```

## `psql` Does **NOT** Support `lock_timeout` in DSN

See `libpq` [Parameter Key Words][2]

```
$ psql "postgres://superuser:testpassword_superuser@localhost:28007/superuser_db?lock_timeout=10ms&sslmode=disable"
psql: error: could not connect to server: invalid URI query parameter: "lock_timeout"
$ psql "postgres://superuser:testpassword_superuser@localhost:28007/superuser_db?sslmode=disable"
psql (12.4, server 10.6)
Type "help" for help.

superuser_db=# \q
$
$
$ psql "user=superuser password=testpassword_superuser host=localhost port=28007 dbname=superuser_db sslmode=disable lock_timeout=10ms"
psql: error: could not connect to server: invalid connection option "lock_timeout"
$ psql "user=superuser password=testpassword_superuser host=localhost port=28007 dbname=superuser_db sslmode=disable"
psql (12.4, server 10.6)
Type "help" for help.

superuser_db=# \q
```

## Clean Up

```
docker rm --force dev-postgres-prevent-deadlock
```

[1]: https://www.postgresql.org/docs/10/errcodes-appendix.html
[2]: https://www.postgresql.org/docs/10/libpq-connect.html#LIBPQ-PARAMKEYWORDS
