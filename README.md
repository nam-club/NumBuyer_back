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
$ http-server -p 3000

http://localhost:3000/index.html
へアクセス
```

### VSCodeでのデバッグ
下記ブログの手順に従う。必要なファイルは既に用意済み。
https://hodalog.com/remote-debug-a-containerized-go-application-using-docker-compose/

## デプロイ
### 構成
コードをS3にアップロードし、CodePipelineで変更を検知、EC2にデプロイ、といった構成。
 
### 初回準備
awsで新規アカウント作成してからコードが自動でサーバに反映されるようになるまでの手順。

1. IAMユーザからアクセスキーID, シークレットアクセスキーを発行し、GitHubのsecretsに設定する。
それぞれ変数名は`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`。
2. GitHubActionsで`Initialize infrastracture`を実行する。
3. `numbuyer-cfn.yaml`のOutputからElastiCacheのホスト名を取得、docker-compose.prd.ymlにセットする。
4. [SSL化対応手順](ssl/SSL化対応手順.md)を参考にEC2をhttps対応にする。

あとはmasterブランチにpushしたのをトリガーに自動でデプロイされていく。

### メンテナンス
- 3ヶ月でSSL証明書の期限が切れるので、ZeroSSLで証明書を再発行、EC2に反映する ※手順後ほど記載