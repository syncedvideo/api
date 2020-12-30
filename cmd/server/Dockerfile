FROM golang:alpine AS builder
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
WORKDIR /build
COPY . .
RUN go mod download

WORKDIR /build/cmd/server
RUN go build -o main .

WORKDIR /dist
RUN cp /build/cmd/server/main .

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /dist /
ENTRYPOINT ["/main"]
