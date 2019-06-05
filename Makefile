PREFIX			?= $(shell pwd)
PKGS 			:= $(shell go list ./... | grep -v /vendor/)
# We don't lint yaml because its forked code and terrible at lint
LINTPKGS        := $(shell go list ./... | grep -v /vendor/ | grep -v "go-sdk/yaml")
SHASUMCMD 		:= $(shell command -v sha1sum || command -v shasum; 2> /dev/null)
TARCMD 			:= $(shell command -v tar || command -v tar; 2> /dev/null)
GIT_REF 		:= $(shell git log --pretty=format:'%h' -n 1)
CURRENT_USER 	:= $(shell whoami)
VERSION 		:= $(shell cat ./VERSION)

# this is to allow local go-sdk/db tests to pass
DB_PORT 		?= 5432
DB_SSLMODE		?= disable

# coverage stuff
CIRCLE_ARTIFACTS 	?= "."
COVERAGE_OUT 		:= "$(CIRCLE_ARTIFACTS)/coverage.html"

GOMAXPROCS := 1

export GOMAXPROCS
export GIT_REF
export VERSION
export DB_SSLMODE

all: format vet profanity test

ci: vet profanity cover-ci

new-install: deps install-all

deps:
	@go get ./...

dev-deps:
	@go get -d github.com/goreleaser/goreleaser

install-all: install-ask install-coverage install-profanity install-reverseproxy install-recover install-semver install-shamir install-template

install-ask:
	@go install github.com/blend/go-sdk/cmd/ask

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

format:
	@echo "$(VERSION)/$(GIT_REF) >> formatting code"
	@go fmt $(PKGS)

vet:
	@echo "$(VERSION)/$(GIT_REF) >> vetting code"
	@go vet $(PKGS)

lint:
	@echo "$(VERSION)/$(GIT_REF) >> linting code"
	@golint $(LINTPKGS)

build:
	@echo "$(VERSION)/$(GIT_REF) >> linting code"
	@docker build . -t go-sdk:$(GIT_REF)
	@docker run go-sdk:$(GIT_REF)

.PHONY: profanity
profanity:
	@echo "$(VERSION)/$(GIT_REF) >> profanity"
	@go run cmd/profanity/main.go --rules PROFANITY_RULES.yml --exclude="cmd/*" --exclude="coverage.html" --exclude="dist/*" --exclude="*/node_modules/*"

test-circleci:
	@echo "$(VERSION)/$(GIT_REF) >> tests"
	@circleci build

test:
	@echo "$(VERSION)/$(GIT_REF) >> tests"
	@go test $(PKGS) -timeout 15s

test-verbose:
	@echo "$(VERSION)/$(GIT_REF) >> tests"
	@go test -v $(PKGS)

cover:
	@echo "$(VERSION)/$(GIT_REF) >> coverage"
	@go run cmd/coverage/main.go --exclude="examples/*"

cover-ci:
	@echo "$(VERSION)/$(GIT_REF) >> coverage"
	@go run cmd/coverage/main.go --keep-coverage-out --covermode=atomic --coverprofile=coverage.txt --exclude="examples/*"

cover-enforce:
	@echo "$(VERSION)/$(GIT_REF) >> coverage"
	@go run cmd/coverage/main.go -enforce --exclude="examples/*"

cover-update:
	@echo "$(VERSION)/$(GIT_REF) >> coverage"
	@go run cmd/coverage/main.go -update --exclude="examples/*"

increment-patch:
	@echo "Current Version $(VERSION)"
	@go run cmd/semver/main.go increment patch ./VERSION > ./NEW_VERSION
	@mv ./NEW_VERSION ./VERSION
	@cat ./VERSION

increment-minor:
	@echo "Current Version $(VERSION)"
	@go run cmd/semver/main.go increment minor ./VERSION > ./NEW_VERSION
	@mv ./NEW_VERSION ./VERSION
	@cat ./VERSION

increment-major:
	@echo "Current Version $(VERSION)"
	@go run cmd/semver/main.go increment major ./VERSION > ./NEW_VERSION
	@mv ./NEW_VERSION ./VERSION
	@cat ./VERSION

clean: clean-dist clean-coverage clean-cache

clean-coverage:
	@echo "Cleaning COVERAGE files"
	@find . -name "COVERAGE" -exec rm {} \;

clean-cache:
	@go clean ./...

clean-dist:
	@rm -rf dist

tag:
	@echo "Tagging v$(VERSION)"
	@git tag -f v$(VERSION)

push-tag:
	@echo "Pushing v$(VERSION) tag to remote"
	@git push -f origin v$(VERSION)

release-all: clean-dist release-ask release-coverage release-job release-profanity release-proxy release-recover release-semver release-shamir release-template

release-ask:
	@goreleaser release -f .goreleaser/ask.yml

release-coverage:
	@goreleaser release -f .goreleaser/coverage.yml

release-job:
	@goreleaser release -f .goreleaser/job.yml

release-profanity:
	@goreleaser release -f .goreleaser/profanity.yml

release-proxy:
	@goreleaser release -f .goreleaser/proxy.yml

release-recover:
	@goreleaser release -f .goreleaser/recover.yml

release-semver:
	@goreleaser release -f .goreleaser/semver.yml

release-shamir:
	@goreleaser release -f .goreleaser/shamir.yml

release-template:
	@goreleaser release -f .goreleaser/template.yml
