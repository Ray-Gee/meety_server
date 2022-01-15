FROM golang:1.13.1 AS builder
WORKDIR /go/src/github.com/Ryuichi-g/meety_server
COPY . .
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
# Build
RUN go build -o app server/main.go

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/Ryuichi-g/meety_server/app /app
EXPOSE 50051
ENTRYPOINT ["/app"]