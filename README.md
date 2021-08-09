# NumBuyerのバックエンドのリポジトリ

## 起動まで

### 事前準備

- Docker, docker-composeのインストール ([参考](https://awesome-linus.com/2019/08/17/mac-docker-install/))

### 各種コマンド
```
起動
$ docker-compose up

停止
$ docker-compose down

キャッシュ無しで再作成
$ docker-compose build --no-cache
```

## サンプル実行手順（後で消す）
```
[事前準備]
以下コマンドを実行しておく。
$ npm install -g http-server

1. Dockerfile で ENTRYPOINT に sample.go を指定する。
2. docker-compose up でサーバを起動しておく。
3. プロジェクトのルートフォルダへ移動し、以下のコマンドを実行する。
$ http-server

http://127.0.0.1:8080/sample_index.html
へアクセス
```