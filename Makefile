VERSION=`git describe --tags`
BUILD=`date +%FT%T%z`
IMAGE_CONFIG=`cat image_config.yml`

LDFLAGS=-ldflags "-w -s -X main.VERSION=${VERSION} -X main.BUILD=${BUILD} -X 'github.com/zdnscloud/zke/types.imageConfig=${IMAGE_CONFIG}'"
GOSRC = $(shell find . -type f -name '*.go')

build: zke

zke: $(GOSRC)
	go build ${LDFLAGS}

clean:
	rm -f zke