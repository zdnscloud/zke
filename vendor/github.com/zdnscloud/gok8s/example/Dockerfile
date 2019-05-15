FROM golang:alpine AS build

RUN mkdir -p /go/src/github.com/zdnscloud/gok8s
COPY . /go/src/github.com/zdnscloud/gok8s

WORKDIR /go/src/github.com/zdnscloud/gok8s/podwatcher
RUN CGO_ENABLED=0 GOOS=linux go build 

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=build /go/src/github.com/zdnscloud/gok8s/podwatcher/podwatcher /usr/local/bin/

ENTRYPOINT ["podwatcher"]
