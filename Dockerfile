ARG GO_VERSION=1.25
ARG CMD_NAME="microservices-api"

FROM golang:${GO_VERSION}-alpine AS build
ARG CMD_NAME

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

RUN echo "Building service: ${CMD_NAME}"

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-s -w" \
  -o /bin/app ./cmd/${CMD_NAME}/main.go


FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
RUN adduser -D -g '' appuser

WORKDIR /app
COPY --from=build /bin/app /app/

RUN chown -R appuser:appuser /app && chmod +x /app/app
USER appuser

ENV PORT=8080
EXPOSE ${PORT}

CMD ["./app"]
