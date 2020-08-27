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
0.000065 ==================================================
0.000096 Configured lock timeout:      10ms
0.000100 Configured context timeout:   600ms
0.000103 Configured transaction sleep: 200ms
0.000124 ==================================================
0.016383 DSN="postgres://superuser:testpassword_superuser@localhost:28007/superuser_db?lock_timeout=10ms&sslmode=disable"
0.016393 ==================================================
0.017589 lock_timeout=10ms
0.260475 ==================================================
0.260484 Error(s):
0.260550 - &pq.Error{Severity:"ERROR", Code:"55P03", Message:"canceling statement due to lock timeout", Detail:"", Hint:"", Position:"", InternalPosition:"", InternalQuery:"", Where:"while updating tuple (0,2) in relation \"might_deadlock\"", Schema:"", Table:"", Column:"", DataTypeName:"", Constraint:"", File:"postgres.c", Line:"2989", Routine:"ProcessInterrupts"}
0.260562 - &pq.Error{Severity:"ERROR", Code:"55P03", Message:"canceling statement due to lock timeout", Detail:"", Hint:"", Position:"", InternalPosition:"", InternalQuery:"", Where:"while updating tuple (0,1) in relation \"might_deadlock\"", Schema:"", Table:"", Column:"", DataTypeName:"", Constraint:"", File:"postgres.c", Line:"2989", Routine:"ProcessInterrupts"}
0.260568 - &errors.errorString{s:"pq: Could not complete operation in a failed transaction"}
0.260572 - &errors.errorString{s:"pq: Could not complete operation in a failed transaction"}
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
0.000057 ==================================================
0.000090 Configured lock timeout:      10s
0.000095 Configured context timeout:   10s
0.000101 Configured transaction sleep: 200ms
0.000128 ==================================================
0.014709 DSN="postgres://superuser:testpassword_superuser@localhost:28007/superuser_db?lock_timeout=10000ms&sslmode=disable"
0.014720 ==================================================
0.017099 lock_timeout=10s
1.254235 ==================================================
1.254245 Error(s):
1.254323 - &pq.Error{Severity:"ERROR", Code:"40P01", Message:"deadlock detected", Detail:"Process 278 waits for ShareLock on transaction 772; blocked by process 277.\nProcess 277 waits for ShareLock on transaction 773; blocked by process 278.", Hint:"See server log for query details.", Position:"", InternalPosition:"", InternalQuery:"", Where:"while updating tuple (0,1) in relation \"might_deadlock\"", Schema:"", Table:"", Column:"", DataTypeName:"", Constraint:"", File:"deadlock.c", Line:"1140", Routine:"DeadLockReport"}
1.254332 - &errors.errorString{s:"pq: Could not complete operation in a failed transaction"}
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
0.000057 ==================================================
0.000086 Configured lock timeout:      10s
0.000090 Configured context timeout:   100ms
0.000103 Configured transaction sleep: 200ms
0.000130 ==================================================
0.014219 DSN="postgres://superuser:testpassword_superuser@localhost:28007/superuser_db?lock_timeout=10000ms&sslmode=disable"
0.014227 ==================================================
0.015685 lock_timeout=10s
0.242769 ==================================================
0.242782 Error(s):
0.242788 - context.deadlineExceededError{}
0.242813 - Context cancel in between queries
0.242817 - context.deadlineExceededError{}
0.242820 - Context cancel in between queries
0.242832 - &errors.errorString{s:"sql: transaction has already been committed or rolled back"}
0.242837 - &errors.errorString{s:"sql: transaction has already been committed or rolled back"}
```

## Cancel "Stuck" Deadlock via Go `context` Cancelation

```
$ DISABLE_LOCK_TIMEOUT=true go run .
0.000051 ==================================================
0.000082 Configured lock timeout:      10s
0.000086 Configured context timeout:   600ms
0.000089 Configured transaction sleep: 200ms
0.000112 ==================================================
0.013133 DSN="postgres://superuser:testpassword_superuser@localhost:28007/superuser_db?lock_timeout=10000ms&sslmode=disable"
0.013142 ==================================================
0.014593 lock_timeout=10s
0.609575 ==================================================
0.609594 Error(s):
0.609649 - &pq.Error{Severity:"ERROR", Code:"57014", Message:"canceling statement due to user request", Detail:"", Hint:"", Position:"", InternalPosition:"", InternalQuery:"", Where:"while updating tuple (0,2) in relation \"might_deadlock\"", Schema:"", Table:"", Column:"", DataTypeName:"", Constraint:"", File:"postgres.c", Line:"3026", Routine:"ProcessInterrupts"}
0.609666 - &pq.Error{Severity:"ERROR", Code:"57014", Message:"canceling statement due to user request", Detail:"", Hint:"", Position:"", InternalPosition:"", InternalQuery:"", Where:"while updating tuple (0,1) in relation \"might_deadlock\"", Schema:"", Table:"", Column:"", DataTypeName:"", Constraint:"", File:"postgres.c", Line:"3026", Routine:"ProcessInterrupts"}
0.609676 - &errors.errorString{s:"sql: transaction has already been committed or rolled back"}
0.609683 - &errors.errorString{s:"sql: transaction has already been committed or rolled back"}
```

From [Appendix A. PostgreSQL Error Codes][1]:

```
Class 57 - Operator Intervention
---------+-----------------------
   57014 | query_canceled
```

## `psql` Does **NOT** Support `lock_timeout` in DSN

```
$ psql "postgres://superuser:testpassword_superuser@localhost:28007/superuser_db?lock_timeout=10ms&sslmode=disable"
psql: error: could not connect to server: invalid URI query parameter: "lock_timeout"
$ psql "postgres://superuser:testpassword_superuser@localhost:28007/superuser_db?sslmode=disable"
psql (12.4, server 10.6)
Type "help" for help.

superuser_db=# \q
```

## Clean Up

```
docker rm --force dev-postgres-prevent-deadlock
```

[1]: https://www.postgresql.org/docs/10/errcodes-appendix.html
