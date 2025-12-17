FROM golang:1.25.5-alpine3.23 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o auth-server ./cmd/server/main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates
RUN mkdir -p /app/internal/config
COPY --from=build /app/auth-server /app/
COPY --from=build /app/internal/config/config.yaml /app/internal/config/
WORKDIR /app
ENTRYPOINT ["./auth-server"]