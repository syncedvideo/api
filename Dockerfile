FROM golang:1.16-alpine as base

RUN apk update && \
    apk upgrade && \
    apk add --no-cache make

# ==========================================

FROM base as ci

ENV CGO_ENABLED=0
COPY . /src
WORKDIR /src

# ==========================================

FROM base as dev

RUN apk update && \
    apk upgrade && \
    apk add --no-cache git 

RUN go get github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest\
    github.com/ramya-rao-a/go-outline \
    github.com/go-delve/delve/cmd/dlv \
    golang.org/x/lint/golint \
    github.com/josharian/impl
RUN GO111MODULE=on go get golang.org/x/tools/gopls@master golang.org/x/tools@master
RUN go install honnef.co/go/tools/cmd/staticcheck@latest
ENV CGO_ENABLED=0

# ==========================================

FROM base AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN mkdir /src/
WORKDIR /src/

COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

COPY . .
WORKDIR /src/cmd/server
RUN go build -o server

# ==========================================

FROM alpine as runtime
COPY --from=builder /src/cmd/server /
CMD ./server
