# EC2で無料https化をする手順

※無料で対応するための暫定手順。3ヶ月で証明書期限が切れてしまい、その度に手動で発効が必要なので、ある程度アクセスが来るようになったら今後はACMなどによって対応すること。

1. ZeroSSLで証明書発行
https://app.zerossl.com/dashboard
- ↑から無料プランで発行していく
- 検証手順は `HTTP File Upload` にする
-  リポジトリの `(Numbuyer_backルート)/ssl/docs` フォルダ以下に指定された検証用ファイルを入れてPUSHし、サーバにデプロイする
-  デプロイが終わったらZeroSSLで検証をする
-  証明書一式zipをダウンロード、展開

2. リポジトリのSSL証明書を更新

下記ディレクトリ以下の証明書を更新しPUSH
```
(Numbuyer_backルート)/resrources/cert
```