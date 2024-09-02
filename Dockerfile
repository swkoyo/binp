# Fetch
FROM golang:1.22.4-alpine3.20 AS fetch-stage
COPY go.mod go.sum /app/
WORKDIR /app
RUN apk add --no-cache gcc musl-dev
RUN go mod download

# Generate Templ
FROM ghcr.io/a-h/templ:v0.2.771 AS generate-templ-stage
COPY --chown=65532:65532 . /app/
WORKDIR /app
RUN ["templ", "generate"]

# Generate Tailwind
FROM node:20.17.0-alpine3.20 AS generate-tailwind-stage
WORKDIR /app
COPY --from=generate-templ-stage /app .
COPY package.json package-lock.json /app/
COPY tailwind.config.js /app/
RUN npm ci
RUN npm run build

# Build
FROM golang:1.22.4-alpine3.20 AS build-stage
RUN apk add --no-cache gcc musl-dev sqlite-dev
COPY --from=generate-tailwind-stage /app /app
WORKDIR /app
ENV CGO_ENABLED=1
ENV GOOS=linux
RUN go build -o /app/bin/api /app/cmd/api/main.go
RUN go build -o /app/bin/cron /app/cmd/cron/main.go

# Final Stage
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata sqlite-libs
RUN adduser -D appuser
WORKDIR /app
COPY --from=build-stage /app/bin/api .
COPY --from=build-stage /app/bin/cron .
COPY --from=build-stage /app/static /app/static
COPY entrypoint.sh .

ENV PORT=8080
ENV GO_ENV=production
ENV DB_PATH=/app/data/db.sqlite

RUN mkdir -p /app/data /app/static/css && \
    chown -R appuser:appuser /app && \
    chmod -R 755 /app/static && \
    chmod +x entrypoint.sh

USER appuser

EXPOSE 8080

ENTRYPOINT ["/bin/sh"]
CMD ["/app/entrypoint.sh"]
