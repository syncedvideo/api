FROM golang:alpine AS builder
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /src
ADD go.mod .
ADD go.sum .
RUN go mod download

ADD . .
WORKDIR /src/cmd/server
RUN go build -o main .

WORKDIR /dist
RUN cp /src/cmd/server/main .

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /dist /
ENTRYPOINT ["/main"]
