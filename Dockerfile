FROM docker.io/library/golang:1.25-alpine AS builder
WORKDIR /app
COPY . /app
ENV CGO_ENABLED=0
ARG version=development
RUN go mod download
RUN go build -trimpath -tags timetzdata \
    -ldflags "-s -w -X github.com/gizmo-ds/misstodon/internal/global.AppVersion=$version" \
    -o misstodon \
    ./cmd/misstodon

FROM gcr.io/distroless/static-debian12:latest
WORKDIR /app
COPY --from=builder /app/misstodon /app/misstodon
COPY --from=builder /app/config_example.toml /app/config.toml
ENTRYPOINT ["/app/misstodon", "start"]
EXPOSE 3000
