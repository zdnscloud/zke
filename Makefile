VERSION=`git describe --tags`
BUILD=`date +%FT%T%z`
IMAGE_INITER_FILE=types/initer.go

LDFLAGS=-ldflags "-w -s -X main.VERSION=${VERSION} -X main.BUILD=${BUILD}"
GOSRC = $(shell find . -type f -name '*.go')

build: zke

zke: $(GOSRC)
	if [ -f $(IMAGE_INITER_FILE) ]; then rm $(IMAGE_INITER_FILE); fi
	go generate types/generate.go
	go build ${LDFLAGS}
	rm -f $(IMAGE_INITER_FILE)

clean:
	rm -f zke