FROM golang:1.17-alpine AS builder

WORKDIR /app
ADD . /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /russ .

FROM alpine:3
RUN apk add --no-cache curl
EXPOSE 8000
COPY --from=builder /russ /russ
RUN chmod +x /russ
ENTRYPOINT ["/russ"]
