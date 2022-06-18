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
3. `numbuyer-cfn.yaml`のOutput中のpublic IPアドレスで、フロントリポジトリの`package.json`にある`REACT_APP_SOCKET_URL`を更新。
3. `numbuyer-cfn.yaml`のOutput中のpublic IPアドレスで、[SSL化対応手順](ssl/SSL化対応手順.md)を参考にEC2をhttps対応にする。

あとはmasterブランチにpushしたのをトリガーに自動でデプロイされていく。

### メンテナンス
- 3ヶ月でSSL証明書の期限が切れるので、ZeroSSLで証明書を再発行、EC2に反映する ※手順後ほど記載

## 設計
### 基本方針
MVC + Service + Repository の構成とする。各層の関係は下記。
```
ユーザの操作 - View(フロント) - Controller - Service - Repository - DB,ネットワーク通信
                                 |           |          |
                                 ------------|-----------
                                           Model       
```
本バックエンドにおける各層とパッケージの対応は下記。
```
Controller: routes配下
Service: services配下
Model:  models配下
Repository: db配下
```

### ゲーム進行
プレイヤーが接続してからゲームを終了するまでのフローは下記。
1. サーバへsocket接続
2. ゲーム作成 or 参加
3. プレイヤーが揃ったらゲーム開始
4. ゲーム終了

全ての制御はroutes配下に定義されたAPIを通して行われるが、3 以降、ゲームの状態を監視し定期的に全プレイヤーに更新状態を連携する必要がある。<br>
その役割はスケジューラ `services/logic/phase_scheduler.go` が担当する。

### DBデータの扱い方、排他制御

■ 排他制御の考え方<br>
本DBはRedisを採用しており、データベースとしてはゲーム用、プレイヤー用で分かれている。

排他制御のためにDBのロックを多用してしまうと下記の問題がある。
- パフォーマンスが落ちる
- コードが複雑になる

そのため、極力排他制御をしないよう、下記の設計にする。

■ ゲーム中以外<br>
- ゲーム用DB
  - ゲームのオーナーに当たるプレイヤーのみSET操作可能
- プレイヤー用DB
  - 各プレイヤーが、自身のハッシュフィールドのみSET操作可能

■ ゲーム中<br>
- ゲーム用DB
  - スケジューラのみSET操作可能
- プレイヤー用DB
  - 各プレイヤーが、自身のハッシュフィールドのみSET操作可能。スケジューラもSET操作が可能のため、mutexによる排他制御を取る。スケジューラはロックされていたらリトライし、プレイヤーはロックされていたらエラーを返却する。