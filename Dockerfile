FROM golang:1.17-alpine AS builder

WORKDIR /app
ADD . /app
RUN go build -o russ .

FROM alpine:3

WORKDIR /app
COPY --from=builder /app/russ /app/russ

ENTRYPOINT ["/app/russ"]