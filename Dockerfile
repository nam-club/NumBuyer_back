FROM golang:1.17
RUN mkdir /go/src/work && \
    go env -w GO111MODULE=on && \
    # ホットリロードの設定
    go install github.com/cosmtrek/air@latest
    
WORKDIR /go/src/work
ADD . /go/src/work

CMD ["air", "-c", ".air.toml"]