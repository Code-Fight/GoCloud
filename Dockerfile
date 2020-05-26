# build golang builder
FROM golang:alpine as builder

RUN apk update

RUN mkdir /app
WORKDIR /app
COPY src/go.mod .
COPY src/go.sum .

RUN go mod download
COPY src .

RUN go build -o /main ./cmd/main.go


# running built service
FROM alpine:3.9

RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

COPY --from=builder /main .
COPY --from=builder /app/web ./web

ENTRYPOINT ["/main"]