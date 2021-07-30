PREFIX				?= $(shell pwd)
PKGS 					:= $(shell go list ./... | grep -v /vendor/)
SHASUMCMD 		:= $(shell command -v sha1sum || command -v shasum; 2> /dev/null)
TARCMD 				:= $(shell command -v tar || command -v tar; 2> /dev/null)
GIT_REF 			:= $(shell git log --pretty=format:'%h' -n 1)
CURRENT_USER 	:= $(shell whoami)
VERSION 			:= $(shell cat ./VERSION)

# this is to allow local go-sdk/db tests to pass
DB_PORT 			?= 5432
DB_SSLMODE		?= disable

export GIT_REF
export VERSION
export DB_SSLMODE

all: ci

ci: deps vet profanity copyright test

new-install: deps dev-deps install-all

list-deps:
	@go list -f '{{ join .Imports "\n" }}' ./... | egrep "github.com|golang.org" | grep -v "github.com/blend/go-sdk" | sort | uniq

deps:
	@go get ./...

dev-deps:
	@go install golang.org/x/lint/golint@latest
	@go install github.com/goreleaser/goreleaser@latest

install-all: install-ask install-copyright install-coverage install-profanity install-reverseproxy install-recover install-semver install-shamir install-template

install-ask:
	@go install github.com/blend/go-sdk/cmd/ask

install-copyright:
	@go install github.com/blend/go-sdk/cmd/copyright

install-coverage:
	@go install github.com/blend/go-sdk/cmd/coverage

install-profanity:
	@go install github.com/blend/go-sdk/cmd/profanity

install-reverseproxy:
	@go install github.com/blend/go-sdk/cmd/reverseproxy

install-recover:
	@go install github.com/blend/go-sdk/cmd/recover

install-semver:
	@go install github.com/blend/go-sdk/cmd/semver

install-shamir:
	@go install github.com/blend/go-sdk/cmd/shamir

install-template:
	@go install github.com/blend/go-sdk/cmd/template

release-binaries:
	@goreleaser ./goreleaser.yml

format:
	@echo "$(VERSION)/$(GIT_REF) >> formatting code"
	@go fmt $(PKGS)

vet:
	@echo "$(VERSION)/$(GIT_REF) >> vetting code"
	@go vet $(PKGS)

lint:
	@echo "$(VERSION)/$(GIT_REF) >> linting code"
	@golint $(PKGS)

generate:
	@echo "$(VERSION)/$(GIT_REF) >> generating code"
	@go generate $(PKGS)

build:
	@echo "$(VERSION)/$(GIT_REF) >> linting code"
	@docker build . -t go-sdk:$(GIT_REF)
	@docker run go-sdk:$(GIT_REF)

.PHONY: profanity
profanity:
	@echo "$(VERSION)/$(GIT_REF) >> profanity"
	@go run cmd/profanity/main.go --rules ".profanity.yml" --exclude-dir="cmd/*" --exclude-file="coverage.html" --exclude-dir="dist/*" --include-file="*.go" --exclude-dir="*/node_modules/*" --exclude-dir="vendor/*" --exclude-dir="examples/*" --verbose

.PHONY: copyright 
copyright:
	@echo "$(VERSION)/$(GIT_REF) >> copyright"
	@go run cmd/copyright/main.go --restrictions-open-source

test:
	@echo "$(VERSION)/$(GIT_REF) >> tests"
	@go test $(PKGS) -timeout 15s

test-race:
	@echo "$(VERSION)/$(GIT_REF) >> tests"
	@go test $(PKGS) -timeout 15s -race

test-verbose:
	@echo "$(VERSION)/$(GIT_REF) >> tests"
	@go test -v $(PKGS) -timeout 15s -race

cover:
	@echo "$(VERSION)/$(GIT_REF) >> coverage"
	@go run cmd/coverage/main.go --exclude="examples/*,cmd/*" --timeout="30s" --race

cover-ci:
	@echo "$(VERSION)/$(GIT_REF) >> coverage"
	@go run cmd/coverage/main.go --keep-coverage-out --covermode=atomic --coverprofile=coverage.txt --exclude="examples/*,cmd/*" --timeout="30s"

clean: clean-dist clean-coverage clean-cache

clean-coverage:
	@echo "Cleaning COVERAGE files"
	@find . -name "COVERAGE" -exec rm {} \;

clean-cache:
	@go clean ./...

clean-dist:
	@rm -rf dist

push: 
	@echo "Tagging $(VERSION)"
	@git add .
	@git commit -am 'Updates from Blend'
	@git tag -f $(VERSION)
	@git push -f origin $(VERSION)
	@git push -f origin HEAD

push-verified: 
	@echo "Tagging $(VERSION)"
	@git add .
	@git commit -S -am 'Updates from Blend'
	@git tag -f $(VERSION)
	@git push -f origin $(VERSION)
	@git push -f origin HEAD
