FROM golang:1.24-alpine AS build
RUN apk add --no-cache alpine-sdk

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o main cmd/api/main.go

FROM alpine:latest AS backend
WORKDIR /app
COPY --from=build /app/main /app/main
EXPOSE ${PORT}
CMD ["./main"]


FROM node:23 AS frontend_builder
WORKDIR /frontend

COPY frontend/package*.json ./
RUN npm install
COPY frontend/. .
RUN npm run build

# FROM node:23-slim AS frontend
# RUN npm install -g serve
# COPY --from=frontend_builder /frontend/dist /app/dist
# EXPOSE 5173
# CMD ["serve", "-s", "/app/dist", "-l", "5173"]

FROM nginx:alpine AS frontend
COPY --from=frontend_builder /frontend/dist /usr/share/nginx/html
COPY nginx-frontend.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
