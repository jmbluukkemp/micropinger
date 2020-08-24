FROM golang:alpine AS builder
RUN apk --no-cache add ca-certificates
WORKDIR $GOPATH/src/mypackage/myapp/
COPY . .
RUN go get -d
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w' -o /usr/bin/micropinger *.go
RUN chmod +x /usr/bin/micropinger

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/bin/micropinger /usr/bin/micropinger
