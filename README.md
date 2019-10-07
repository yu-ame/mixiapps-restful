# mixiapps-restful

mixiアプリ　Restful API 署名検証用アプリ

## 必要な要件

dockerがインストールされたサーバー

## 構築

```
git clone https://github.com/yu-ame/mixiapps-restful.git
cd mixiapps-restful

#defaultのファイルをコピーして起動ポートやURLを設定してください。
cp configs/config_default.json configs/config.json 
vi configs/config.json

#ビルド
docker build -f build/dockerfiles/Dockerfile . -t mixiapps-restful

#実行
docker run --rm -p [指定したポート]:[指定したポート] mixiapps-restful
```

外部サーバへのイメージのコピーは下記のようなコマンドが楽です

```
docker save mixiapps-restful > /tmp/mixiapps-restful.tar
scp /tmp/mixiapps-restful.tar user@host:/tmp/
ssh user@host docker load < /tmp/mixiapps-restful.tar
```
