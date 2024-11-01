# signature_downloader

公式ミラーサイトから署名ファイルをダウンロードし、S3バケットにアップロードするためのツールです。
/fleshclam/conf.d/clamav.conf.templateでプライベートS3のURLを設定する場合に使用します。
コンテナ起動後、公式ミラーからsignatureをダウンロードし、指定のS3バケットにアップロードして終了します

## 使用方法
- ECSタスクなどで定期期実行することを想定しています
- ECSやS3での静的ウェブサイトホスティングなどのインフラ設定についてはここでは割愛します

### 1. Imageのビルド
以下コマンドにて、Dockerイメージをビルドし、ECRなどにプッシュします。
```bash
$ docker build  .

# M1以降のMacなど, arm64環境の場合は以下のようにビルドします
$ docker build --platform linux/amd64 .
```

### 2. タスク定期の設定
タスク定義を作成します。タスク定義には以下の環境変数設定を含めます。

- MIRROR_SITE_S3_BUCKET_NAME: プライベートミラーサイトとして使用するS3のバケット名

また、タスクのロールには、上記バケットへのs3:PutObject権限を持つIAMロールを指定します。

### 3. タスクの実行
EventBridgeなどで定期実行するように設定します。

## ローカルでの動作確認
./docker-compose.ymlに適宜環境変数を設定し、以下コマンドで実行します。
```bash
$ docker-compose up
```