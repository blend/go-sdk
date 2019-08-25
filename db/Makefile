all: test

test:
	@DB_SSLMODE=disable go test ./...

bench:
	@DB_SSLMODE=disable go test -bench=.
