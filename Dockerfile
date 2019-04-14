FROM golang:1.12.4 as builder 

# Build
ENV GO111MODULE=on

WORKDIR /app

COPY . .
RUN go mod download

RUN cd /app/cmd/fusio/; GOOS=linux GOARCH=amd64 go build
RUN cp /app/cmd/fusio/fusio /app/fusio


# Image
FROM debian:9
RUN mkdir /etc/fusio /var/log/fusio

COPY --from=builder /app/fusio /app/

EXPOSE 8080
ENTRYPOINT ["/app/fusio", "-c", "/etc/fusio/config.yaml"]

