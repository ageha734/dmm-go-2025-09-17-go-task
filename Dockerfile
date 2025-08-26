FROM golang:1.24.5-alpine3.22 AS builder

WORKDIR /app

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=bind,source=.,target=. \
    CGO_ENABLED=0 GOOS=linux go build -o /bin/app cmd/main.go

FROM gcr.io/distroless/base-debian12:nonroot

# RUN apk --no-cache add ca-certificates

COPY --from=builder /bin/app /bin/app

ENTRYPOINT ["/bin/app"]
