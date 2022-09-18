# syntax=docker/dockerfile:1

FROM golang:1.19-alpine as builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN go build -o ates ./cmd/main.go

## Deploy
FROM alpine:3.16

WORKDIR /app
COPY --from=builder /app/ates /app/ates
EXPOSE 8080
ENTRYPOINT ["/app/ates"]
CMD ["/app/ates"]