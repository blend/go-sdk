all: static

deps:
	@go get -u github.com/kataras/bindata/cmd/bindata

build:
	@bindata -o static.go -pkg=jobkit -md5checksum static/...
