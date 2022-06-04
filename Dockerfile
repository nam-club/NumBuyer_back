FROM golang:1.17
RUN mkdir /go/src/work && \
    go env -w GO111MODULE=on && \
    # 以下、ローカルでのデバッグの設定
    go install golang.org/x/tools/gopls@latest && \
    go install github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest && \
    go install github.com/ramya-rao-a/go-outline@latest && \
    go install github.com/cweill/gotests/gotests@latest && \
    go install github.com/fatih/gomodifytags@latest && \
    go install github.com/josharian/impl@latest && \
    go install github.com/haya14busa/goplay/cmd/goplay@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest && \
    go install github.com/go-delve/delve/cmd/dlv@master && \
    go install honnef.co/go/tools/cmd/staticcheck@latest && \
    go install golang.org/x/tools/gopls@latest && \
    go install github.com/golang/vscode-go@latest && \
    # ホットリロードの設定
    go install -u github.com/cosmtrek/air@latest
    
WORKDIR /go/src/work
ADD . /go/src/work

CMD ["air", "-c", ".air.toml"]