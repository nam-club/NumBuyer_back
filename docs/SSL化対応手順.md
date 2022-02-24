# EC2で無料https化をする手順

※無料で対応するための暫定手順。3ヶ月で証明書期限が切れてしまい、その度に手動で発効が必要なので、ある程度アクセスが来るようになったら今後はACMなどによって対応すること。

1. nginxをEC2にインストール
```
sudo amazon-linux-extras install nginx1 -y
sudo cp -a /etc/nginx/nginx.conf /etc/nginx/nginx.conf.back
sudo systemctl start nginx
sudo systemctl enable nginx
```
2. ZeroSSLで証明書発行

https://app.zerossl.com/dashboard

- ↑から無料プランで発行していく
- 検証手順は `HTTP File Upload` にする

- 検証用のtxtファイルを置くフォルダ作成
```
sudo mkdir -p /usr/share/nginx/html/.well-known/pki-validation/
```
- 上記フォルダ内に指定されたファイルを用意する
- ZeroSSLで検証をする
- 証明書一式zipをダウンロード、展開

3. リポジトリのSSL証明書を更新
下記ディレクトリ以下の証明書を更新しpush
```
(Numbuyer_backルート)/resrources/cert
```

4. nginxはもう使わないので停止する
```
sudo service nginx stop
```