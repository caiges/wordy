FROM golang:alpine as dev
ENV GOCACHE /home/dev/.gocache
ENV GOPATH /home/dev/go
ENV GOPRIVATE=""
ENV PATH="${PATH}:${GOPATH}/bin"
ARG uid=1010
RUN adduser -D -u $uid dev
RUN mkdir -p /home/dev/.gocache
RUN apk add --no-cache \
  curl \
  git \
  openssh \
  alpine-sdk
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.30.0
RUN chown -f -R dev:dev $GOCACHE $GOPATH
WORKDIR /home/dev
