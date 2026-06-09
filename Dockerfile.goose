FROM golang:1.26.4-alpine AS build
RUN apk add --no-cache alpine-sdk

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o main cmd/goose/main.go

FROM alpine:3.22 AS migrations
WORKDIR /app
RUN addgroup -S app && adduser -S -G app app
COPY --from=build --chown=app:app /app/main /app/main
USER app
HEALTHCHECK NONE
CMD ["./main"]
