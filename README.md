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

## ローカルでのデバッグ手順

### フロントテスト用のサーバ起動
```
[事前準備]
以下コマンドを実行しておく。
$ npm install -g http-server

1. docker-compose up でサーバを起動しておく。
2. プロジェクトのルートフォルダへ移動し、以下のコマンドを実行する。
$ http-server

http://127.0.0.1:8080/index.html
へアクセス
```

### VSCodeでのデバッグ
下記ブログの手順に従う。必要なファイルは既に用意済み。
https://hodalog.com/remote-debug-a-containerized-go-application-using-docker-compose/

## デプロイ
### 構成
コードをS3にアップロードし、CodePipelineで変更を検知、EC2にデプロイ、といった構成。
 
### 手順
1. IAMユーザからアクセスキーID, シークレットアクセスキーを発行し、GitHubのsecretsに設定する。
それぞれ変数名は`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`。
2. GitHubActionsで`Initialize infrastracture`を実行する。

あとはmasterブランチにpushしたのをトリガーに自動でデプロイされていく。