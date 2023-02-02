fmtcheck:
	"$(CURDIR)/scripts/gofmtcheck.sh"

importscheck:
	"$(CURDIR)/scripts/goimportscheck.sh"

staticcheck:
	"$(CURDIR)/scripts/staticcheck.sh"

exhaustive:
	"$(CURDIR)/scripts/exhaustive.sh"

build-linux:
	mkdir -p _out
	GOOS=linux GOARCH=amd64 go build -v -o ./_out/unity-control-plane ./cmd/web

build: build-linux

test:
	mkdir -p _out
	go test -p 1 -v ./...

.PHONY: fmtcheck importscheck staticcheck build all

all: fmtcheck importscheck staticcheck build
