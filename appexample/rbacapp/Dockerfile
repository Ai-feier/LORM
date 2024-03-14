FROM golang:1.20 as builder

WORKDIR /app

copy . .

RUN go env -w GOPROXY=https://goproxy.cn,direct && go mod tidy && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go

RUN chmod +x ./main

FROM alpine:3.19

RUN apk add --update bash

WORKDIR /app

COPY --from=builder /app/main .

CMD ["/app/main"]

