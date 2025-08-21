FROM golang:1.24-alpine AS builder

WORKDIR /srv

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /srv

COPY --from=builder /srv/app .

EXPOSE 8080

CMD ["./app"]
