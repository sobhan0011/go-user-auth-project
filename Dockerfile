FROM golang:1.23-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

FROM alpine:latest
RUN adduser -D -s /bin/sh appuser
WORKDIR /
COPY --from=build /app/server /server
COPY --from=build /app/internal/infra/db/postgres/migrations /migrations
COPY --from=build /app/openapi.yaml /openapi.yaml
ENV HTTP_PORT=8080
EXPOSE 8080
USER appuser
ENTRYPOINT ["/server"]


