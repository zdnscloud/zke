VERSION=`git describe --tags`
BUILD=`date +%FT%T%z`

LDFLAGS=-ldflags "-w -s -X main.VERSION=${VERSION} -X main.BUILD=${BUILD}"
GOSRC = $(shell find . -type f -name '*.go')

build: zke

zke: $(GOSRC)
	go generate types/generate.go
	go build ${LDFLAGS}
	rm -f types/initer.go

clean:
	rm -f zke