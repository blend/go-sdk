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

## Intentional Contention

In order to introduce a deadlock, we borrow an example from
[When Postgres blocks: 7 tips for dealing with locks][6].

In the first transaction we update "hello" rows following by "world" rows

```sql
BEGIN;
UPDATE might_deadlock SET counter = counter + 1 WHERE key = 'hello';
-- Sleep for configured transaction sleep
UPDATE might_deadlock SET counter = counter + 1 WHERE key = 'world';
COMMIT;
```

and in the second transaction we update the rows in the opposite order

```sql
BEGIN;
UPDATE might_deadlock SET counter = counter + 1 WHERE key = 'world';
-- Sleep for configured transaction sleep
UPDATE might_deadlock SET counter = counter + 1 WHERE key = 'hello';
COMMIT;
```

## Let `postgres` Cancel Via `lock_timeout`

```
$ go run .
0.000055 ==================================================
0.000085 Configured lock timeout:      10ms
0.000089 Configured context timeout:   600ms
0.000091 Configured transaction sleep: 200ms
0.000114 ==================================================
0.014372 DSN="postgres://superuser:testpassword_superuser@localhost:28007/superuser_db?connect_timeout=5&lock_timeout=10ms&sslmode=disable"
0.014381 ==================================================
0.015569 lock_timeout=10ms
0.026958 ==================================================
0.026981 Starting transactions
0.036223 Transactions opened
0.261793 ***
0.261803 Error(s):
0.261852 - &pq.Error{Severity:"ERROR", Code:"55P03", Message:"canceling statement due to lock timeout", Detail:"", Hint:"", Position:"", InternalPosition:"", InternalQuery:"", Where:"while updating tuple (0,2) in relation \"might_deadlock\"", Schema:"", Table:"", Column:"", DataTypeName:"", Constraint:"", File:"postgres.c", Line:"2989", Routine:"ProcessInterrupts"}
0.261862 - &pq.Error{Severity:"ERROR", Code:"55P03", Message:"canceling statement due to lock timeout", Detail:"", Hint:"", Position:"", InternalPosition:"", InternalQuery:"", Where:"while updating tuple (0,1) in relation \"might_deadlock\"", Schema:"", Table:"", Column:"", DataTypeName:"", Constraint:"", File:"postgres.c", Line:"2989", Routine:"ProcessInterrupts"}
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
0.011839 DSN="postgres://superuser:testpassword_superuser@localhost:28007/superuser_db?connect_timeout=5&lock_timeout=10000ms&sslmode=disable"
0.011850 ==================================================
0.013332 lock_timeout=10s
0.022643 ==================================================
0.022659 Starting transactions
0.030515 Transactions opened
1.245005 ***
1.245016 Error(s):
1.245053 - &pq.Error{Severity:"ERROR", Code:"40P01", Message:"deadlock detected", Detail:"Process 347 waits for ShareLock on transaction 845; blocked by process 346.\nProcess 346 waits for ShareLock on transaction 846; blocked by process 347.", Hint:"See server log for query details.", Position:"", InternalPosition:"", InternalQuery:"", Where:"while updating tuple (0,1) in relation \"might_deadlock\"", Schema:"", Table:"", Column:"", DataTypeName:"", Constraint:"", File:"deadlock.c", Line:"1140", Routine:"DeadLockReport"}
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
0.013163 DSN="postgres://superuser:testpassword_superuser@localhost:28007/superuser_db?connect_timeout=5&lock_timeout=10000ms&sslmode=disable"
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
```

## Cancel "Stuck" Deadlock via Go `context` Cancelation

```
$ DISABLE_LOCK_TIMEOUT=true go run .
0.000053 ==================================================
0.000084 Configured lock timeout:      10s
0.000088 Configured context timeout:   600ms
0.000091 Configured transaction sleep: 200ms
0.000113 ==================================================
0.014431 DSN="postgres://superuser:testpassword_superuser@localhost:28007/superuser_db?connect_timeout=5&lock_timeout=10000ms&sslmode=disable"
0.014442 ==================================================
0.016239 lock_timeout=10s
0.026890 ==================================================
0.026903 Starting transactions
0.036405 Transactions opened
0.612462 ***
0.612474 Error(s):
0.612539 - &pq.Error{Severity:"ERROR", Code:"57014", Message:"canceling statement due to user request", Detail:"", Hint:"", Position:"", InternalPosition:"", InternalQuery:"", Where:"while updating tuple (0,2) in relation \"might_deadlock\"", Schema:"", Table:"", Column:"", DataTypeName:"", Constraint:"", File:"postgres.c", Line:"3026", Routine:"ProcessInterrupts"}
0.612556 - &pq.Error{Severity:"ERROR", Code:"57014", Message:"canceling statement due to user request", Detail:"", Hint:"", Position:"", InternalPosition:"", InternalQuery:"", Where:"while updating tuple (0,1) in relation \"might_deadlock\"", Schema:"", Table:"", Column:"", DataTypeName:"", Constraint:"", File:"postgres.c", Line:"3026", Routine:"ProcessInterrupts"}
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
$ psql "postgres://superuser:testpassword_superuser@localhost:28007/superuser_db?connect_timeout=5&lock_timeout=10ms&sslmode=disable"
psql: error: could not connect to server: invalid URI query parameter: "lock_timeout"
$
$
$ psql "postgres://superuser:testpassword_superuser@localhost:28007/superuser_db?connect_timeout=5&sslmode=disable"
...
superuser_db=# SHOW lock_timeout;
 lock_timeout
--------------
 0
(1 row)

superuser_db=# \q
$
$
$ PGOPTIONS="-c lock_timeout=4500ms" psql "postgres://superuser:testpassword_superuser@localhost:28007/superuser_db?connect_timeout=5&sslmode=disable"
...
superuser_db=# SHOW lock_timeout;
 lock_timeout
--------------
 4500ms
(1 row)

superuser_db=# \q
```

Instead `github.com/lib/pq` parses all query parameters when [reading a DSN][3]
(data source name) and then passes all non-driver setting through as key-value
pairs when [forming a startup packet][4]. The designated driver settings
[are][5]:

- `host`
- `port`
- `password`
- `sslmode`
- `sslcert`
- `sslkey`
- `sslrootcert`
- `fallback_application_name`
- `connect_timeout`
- `disable_prepared_binary_result`
- `binary_parameters`
- `krbsrvname`
- `krbspn`

From the [Start-up][7] section of the documentation for
"Frontend/Backend Protocol > Message Flow", we see

> To begin a session, a frontend opens a connection to the server and sends a
> startup message ... (Optionally, the startup message can include additional
> settings for run-time parameters.)

This is why using `PGOPTIONS="-c {key}={value}"` is required for [setting][8]
named run-time parameters. It's also worth noting that `github.com/lib/pq`
is [`PGOPTIONS`-aware][9].

## Clean Up

```
docker rm --force dev-postgres-prevent-deadlock
```

[1]: https://www.postgresql.org/docs/10/errcodes-appendix.html
[2]: https://www.postgresql.org/docs/10/libpq-connect.html#LIBPQ-PARAMKEYWORDS
[3]: https://github.com/lib/pq/blob/v1.8.0/connector.go#L67-L69
[4]: https://github.com/lib/pq/blob/v1.8.0/conn.go#L1093-L1105
[5]: https://github.com/lib/pq/blob/v1.8.0/conn.go#L1058-L1084
[6]: https://www.citusdata.com/blog/2018/02/22/seven-tips-for-dealing-with-postgres-locks/
[7]: https://www.postgresql.org/docs/10/protocol-flow.html#id-1.10.5.7.3
[8]: https://www.postgresql.org/docs/10/app-postgres.html#id-1.9.5.13.6.3
[9]: https://github.com/lib/pq/blob/v1.8.0/conn.go#L1945-L1946
