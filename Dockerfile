# Build MJML email templates
FROM node:23 AS mjml_builder
WORKDIR /app
COPY . .

RUN npm i -g mjml
RUN mjml internal/email/emails/*.mjml -o internal/email/emails/

# Build API binary
FROM golang:1.26.4-alpine AS build
RUN apk add --no-cache alpine-sdk

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY --from=mjml_builder /app/. .

RUN CGO_ENABLED=1 GOOS=linux go build -o main cmd/api/main.go

FROM alpine:3.22 AS backend
WORKDIR /app
RUN addgroup -S app && adduser -S -G app app
COPY --from=build --chown=app:app /app/main /app/main
USER app
EXPOSE ${PORT}
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget --quiet --tries=1 --spider "http://127.0.0.1:${PORT}/live" || exit 1
CMD ["./main"]

FROM node:23 AS frontend_builder
WORKDIR /frontend

COPY frontend/package*.json ./
RUN npm install
COPY frontend/. .
RUN npm run build

FROM nginx:alpine AS frontend
COPY --from=frontend_builder /frontend/dist /usr/share/nginx/html
COPY nginx-frontend.conf /etc/nginx/conf.d/default.conf
RUN chown -R nginx:nginx /run /var/cache/nginx /usr/share/nginx/html
USER nginx
EXPOSE 80
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://127.0.0.1/ || exit 1
