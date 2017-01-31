FROM golang:1.7-alpine

RUN apk --update add git #&& \
    # go get \
    #   github.com/golang/dep

RUN mkdir -p /go/src/github.com/minodisk/go-learn-sql
WORKDIR /go/src/github.com/minodisk/go-learn-sql
COPY . .

CMD ls -al vendor/github.com/go-sql-driver/mysql && go run main.go
