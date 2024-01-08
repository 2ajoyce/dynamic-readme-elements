FROM golang:1.21-alpine AS builder
WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /usr/local/bin/app

FROM alpine:latest
COPY --from=builder /usr/local/bin/app /app
ENTRYPOINT ["/app"]