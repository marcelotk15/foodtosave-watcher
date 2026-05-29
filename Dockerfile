FROM golang:1.24-alpine AS builder

ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build \
  -ldflags="-s -w" \
  -o /app/bin/foodtosave-watcher \
  ./cmd/foodtosave-watcher

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /app/bin/foodtosave-watcher /app/foodtosave-watcher

USER nonroot:nonroot

ENTRYPOINT ["/app/foodtosave-watcher"]
CMD ["--config", "/app/config.yaml"]
