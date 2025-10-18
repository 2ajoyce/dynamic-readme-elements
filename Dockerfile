FROM golang:1.25 AS builder
WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY v0/go.mod v0/go.sum ./v0/
WORKDIR /usr/src/app/v0
RUN go mod download && go mod verify
WORKDIR /usr/src/app
COPY . .
WORKDIR /usr/src/app/v0
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install .

FROM alpine:latest
COPY --from=builder /go/bin/v0 /app
ENV GIN_MODE=release
ENTRYPOINT ["/app"]