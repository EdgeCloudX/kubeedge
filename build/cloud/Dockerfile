FROM golang:1.16-alpine3.13 AS builder

ARG GO_LDFLAGS

COPY . /go/src/github.com/kubeedge/kubeedge

ENV GOPROXY=https://goproxy.cn
WORKDIR /go/src/github.com/kubeedge/kubeedge
RUN export GO111MODULE=on \
    &&  go mod tidy \
    &&  go mod vendor

RUN CGO_ENABLED=0 GO111MODULE=off go build -v -o /usr/local/bin/cloudcore -ldflags "$GO_LDFLAGS -w -s" \
    ./cloud/cmd/cloudcore


FROM alpine:3.13

COPY --from=builder /usr/local/bin/cloudcore /usr/local/bin/cloudcore

RUN apk add --update-cache \
    iptables \
    && rm -rf /var/cache/apk/*

ENTRYPOINT ["cloudcore"]
