# stage 1: build src code to binary
# 获取对应版本号 https://hub.docker.com/_/golang
FROM golang:1.18-alpine3.15 as builder

# Set up dependencies
ENV PACKAGES make gcc git libc-dev linux-headers bash

COPY  . $GOPATH/src
WORKDIR $GOPATH/src

# Install minimum necessary dependencies, build binary
# RUN apk add --no-cache $PACKAGES && make all
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
apk add --no-cache $PACKAGES && make all

FROM alpine:3.15

COPY --from=builder /go/src/deal-tx-field /usr/local/bin/

CMD ["deal-tx-field", "start"]
