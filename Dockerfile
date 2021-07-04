FROM golang:latest
RUN mkdir /go/src/work && \
    go env -w GO111MODULE=on && \
    go get github.com/googollee/go-socket.io && \
    go get github.com/gomodule/redigo && \
    go get github.com/gin-gonic/gin
WORKDIR /go/src/work
ADD . /go/src/work

# サンプル実行の時は sample.go を指定する
ENTRYPOINT go run main.go