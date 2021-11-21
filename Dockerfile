FROM golang:latest
RUN mkdir /go/src/work && \
    go env -w GO111MODULE=on && \
    # 以下、ローカルでのデバッグ用リポジトリ
    go get golang.org/x/tools/gopls@latest && \
    go get github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest && \
    go get github.com/ramya-rao-a/go-outline@latest && \
    go get github.com/cweill/gotests/gotests@latest && \
    go get github.com/fatih/gomodifytags@latest && \
    go get github.com/josharian/impl@latest && \
    go get github.com/haya14busa/goplay/cmd/goplay@latest && \
    go get github.com/go-delve/delve/cmd/dlv@latest && \
    go get github.com/go-delve/delve/cmd/dlv@master && \
    go get honnef.co/go/tools/cmd/staticcheck@latest && \
    go get golang.org/x/tools/gopls@latest && \
    go get github.com/golang/vscode-go
    
WORKDIR /go/src/work
ADD . /go/src/work

# ENTRYPOINT go run main.go