FROM golang:1.14.0-buster as builder
WORKDIR /go/src/app
ADD src src
RUN (cd src && go mod download)
RUN (cd src && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/src/app/main main.go)

FROM alpine:latest
ENV PORT=80
ENV GIN_MODE=release
RUN apk add --no-cache ca-certificates
EXPOSE 80
WORKDIR /root
COPY --from=builder /go/src/app/main .
CMD ["sh", "-c", "./main"]