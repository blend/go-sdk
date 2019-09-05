PREFIX			?= $(shell pwd)
PKGS 			:= $(shell go list ./... | grep -v /vendor/)
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
BUILD_NUMBER 		?= ${GIT_REF}

export GIT_REF
export BUILD_NUMBER
export VERSION
export DB_SSLMODE

all: format vet profanity test

ci: vet profanity cover

new-install: deps install

deps:
	@go get github.com/lib/pq
	@go get -u ./...

dev-deps:
	@go get -d github.com/goreleaser/goreleaser

install: install-ask install-coverage install-profanity install-proxy install-recover install-semver install-shamir install-template

install-ask:
	@go install github.com/blend/go-sdk/cmd/ask

install-coverage:
	@go install github.com/blend/go-sdk/cmd/coverage

install-profanity:
	@go install github.com/blend/go-sdk/cmd/profanity

install-proxy:
	@go install github.com/blend/go-sdk/cmd/proxy

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
	@golint $(PKGS)

build:
	@echo "$(VERSION)/$(GIT_REF) >> linting code"
	@docker build . -t go-sdk:$(GIT_REF)
	@docker run go-sdk:$(GIT_REF)

.PHONY: profanity
profanity:
	@echo "$(VERSION)/$(GIT_REF) >> profanity"
	@go run cmd/profanity/main.go -rules PROFANITY --exclude="cmd/*,coverage.html,dist/*"

test-circleci:
	@echo "$(VERSION)/$(GIT_REF) >> tests"
	@circleci build

test:
	@echo "$(VERSION)/$(GIT_REF) >> tests"
	@go test $(PKGS) -timeout 15s

test-docker:
	@echo "$(VERSION)/$(GIT_REF) >> tests (docker)"
	@bash ./_bin/run_tests docker-compose.yml

test-verbose:
	@echo "$(VERSION)/$(GIT_REF) >> tests"
	@go test -v $(PKGS)

cover:
	@echo "$(VERSION)/$(GIT_REF) >> coverage"
	@go run cmd/coverage/main.go

cover-enforce:
	@echo "$(VERSION)/$(GIT_REF) >> coverage"
	@go run cmd/coverage/main.go -enforce

cover-update:
	@echo "$(VERSION)/$(GIT_REF) >> coverage"
	@go run cmd/coverage/main.go -update

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

release: clean-dist tag push-tag release-ask release-coverage release-profanity release-proxy release-recover release-semver release-template

tag:
	@echo "Tagging v$(VERSION)"
	@git tag -f v$(VERSION)

push-tag:
	@echo "Pushing v$(VERSION) tag to remote"
	@git push -f origin v$(VERSION)

release-ask:
	@goreleaser release -f .goreleaser/ask.yml

release-coverage:
	@goreleaser release -f .goreleaser/coverage.yml

release-profanity:
	@goreleaser release -f .goreleaser/profanity.yml

release-proxy:
	@goreleaser release -f .goreleaser/proxy.yml

release-recover:
	@goreleaser release -f .goreleaser/recover.yml

release-semver:
	@goreleaser release -f .goreleaser/semver.yml

release-template:
	@goreleaser release -f .goreleaser/template.yml
