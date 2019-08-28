all: test

test:
	@DB_SSLMODE=disable go test ./...

cover:
	@DB_SSLMODE=disable go test -cover

bench:
	@DB_SSLMODE=disable go test -bench=.
