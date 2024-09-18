# go-clamav

go-clamavは、ClamAVを使用してS3バケットにアップロードされたファイルをスキャンするためのソフトウェアです。


## 構成

### 全体構成

S3バケットにファイルがアップロードされると、S3イベント通知がSQSに送信されます。go-clamavはSQSからメッセージを受信し、 ファイルをスキャンします。
スキャン結果はDynamoDBに保存されます。

ClamAVが参照するシグネチャデータベースは、サイドカーコンテナとして起動しているfreshclamコンテナとボリュームを共有することで、定期的に更新されます。


### ローカルでの開発

compose.yamlにある通り、以下のコンテナ構成で開発ができるようになっています。

1. minio
    + [MinIO社](https://min.io/)が開発するS3互換のオブジェクトストレージです
    + [minio/minio](https://github.com/minio/minio)にてコードは公開されています
2. ElasticMQ
    + [SoftwareMill社](https://softwaremill.com/)が開発するAmazon SQS互換のインタフェースを持つインメモリメッセージキューソフトウェアです
    + [softwaremill/elasticmq](https://github.com/softwaremill/elasticmq)にてコードは公開されています。
3. DynamoDB Local + DynamoDB Admin
    + DynamoDB Localは、Amazon Web Serviceが提供するDynamoDBを使ったアプリケーションの開発とテストをローカルで進められるようにするためのツールです。
    + 詳細については、[公式ドキュメント](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBLocal.html)を参照してください。

compose.yamlに記載されているほかコンテナは、上記コンテナの初期設定を行うためのサイドカーコンテナです。

### ローカルでの起動手順

docker composeでgo-clamavの動作に必要な関連コンテナが一通り起動します。

```bash
$ docker compose up -d
```

起動タイミング次第では、DynamoDBのテーブルなどが作成されていない場合があります。
その場合は上記コマンドを再度実行してください。

### 環境変数

go-clamavを動かすためには、`config/config.go`に記載されている環境変数を設定する必要があります。

+ QUEUE_URL
  + SQSキューのURL、ローカル開発の場合は、ElasticMQのキューURLを指定します。
  + compose.yamlを利用しているの場合は、`http://localhost:9324/000000000000/queue1` を指定します。
+ Region
  + 各種リソースが配置されているリージョンを指定します。
  + デフォルトでは、東京リージョン(`ap-northeast-1`)を指定します。
+ MAX_NUMBER_OF_MESSAGES
  + 1回のリクエストで取得するSQSキュー内メッセージの最大数を指定します。
  + デフォルトでは1となっています。
+ VISIBILITY_TIMEOUT
  + SQSキュー内のメッセージの可視性タイムアウトを指定します。
  + デフォルトでは30秒となっています。
  + 非常に大きなファイルをスキャンする場合は、この値を大きくすることを検討してください。
+ WAIT_TIME_SECONDS
  + ロングポーリングを行う際の待機時間を指定します。
  + デフォルトでは20秒となっています。
+ BaseUrl
  + SQSキューにアクセスする際のエンドポイントを指定します。
  + ローカル開発の際に、ElasticMQのエンドポイントを指定するための環境変数です。
  + compose.yamlを利用しているの場合は、`http://localhost:9324` を指定します。
+ DYNAMODB_TABLE_NAME
  + スキャン結果を保存するDynamoDBテーブル名を指定します。
  + デフォルトでは、`ScanResults` となっています。 
+ DYNAMODB_TABLE_INFECTED
  + ウイルス検出されたファイルの情報を保存するDynamoDBテーブル名を指定します。
  + デフォルトでは、`InfectedScanResults` となっています。
+ DynamoDBBaseUrl
  + DynamoDBにアクセスする際のエンドポイントを指定します。
  + ローカル開発の際に、DynamoDB Localのエンドポイントを指定するための環境変数です。
  + compose.yamlを利用しているの場合は、`http://localhost:8000` を指定します。
+ CLAMD_HOST
  + ClamAVのホスト名を指定します。
  + デフォルトでは`localhost`となっています。compose.yamlを利用している場合の指定もデフォルトままで問題ありません。
  + 本番環境では、ClamAVが稼働しているホスト名を指定してください。ECSで単一のタスク定義の中でClamAVとgo-clamavを動かす場合は、`localhost`を指定してください。
+ CLAMD_PORT
  + ClamAVのポート番号を指定します。
  + デフォルトでは`3310`となっています。compose.yamlを利用している場合の指定もデフォルトままで問題ありません。

ほか`config/config.go`に記載されていない環境変数としては以下が必要です。

+ AWS_ACCESS_KEY_ID
  + AWSのアクセスキーを指定します。
  + compose.yamlを使っている場合は、`minio` を指定します。
+ AWS_SECRET_ACCESS_KEY
  + AWSのシークレットアクセスキーを指定します。 
  + compose.yamlを使っている場合は、`miniominio` を指定します。

なお、`AWS_ACCESS_KEY_ID`および`AWS_SECRET_ACCESS_KEY`は、compose.yaml利用時には、minioのブラウザUI(`http://localhost:9001`)にアクセスする際のID/PASSの情報となります。

## compose.yaml利用時のキューの操作

あらかじめminioのテストバケットにファイルをアップロードしておくと、go-clamavの動作確認が可能です。
なお、minioへのファイルアップロードイベントのElasticMQへの通知設定ができていないため、AWS CLIを使ってSQSにメッセージを送信する必要があります。

以降で上げているコマンドを正しく実行するためにはminioのtestバケットに以下のファイルをアップロードする必要があります。

1. eicar.txt
    + [EICARテストファイル](https://ja.wikipedia.org/wiki/EICAR%E3%83%86%E3%82%B9%E3%83%88%E3%83%95%E3%82%A1%E3%82%A4%E3%83%AB)です。
    + ウイルススキャンの動作確認でウイルス感染時の挙動をエミュレートする際に利用します。
2. test.txt
    + 正常なテキストファイルです。中身は何でも構いません。

以降は、ElasticMQを使ってキューを操作するためのコマンドを記載しています。

<details>
  <summary>eicar.txtのアップロードイベントテスト用コマンド</summary>


```
aws sqs send-message --queue-url http://localhost:9324/000000000000/queue1 \
 --endpoint-url http://localhost:9324 \
 --message-body '{"Records": [{
      "eventVersion": "2.1",
      "eventSource": "aws:s3",
      "awsRegion": "ap-northeast-1",
      "eventTime": "2022-08-07T14:33:59.870Z",
      "eventName": "ObjectCreated:Put",
      "userIdentity": {
        "principalId": "AWS:AIDAVMRY2N7OKTN33RYNV"
      },
      "requestParameters": {
        "sourceIPAddress": "60.95.0.122"
      },
      "responseElements": {
        "x-amz-request-id": "Q73VJ1CPJ64CKJQ0",
        "x-amz-id-2": "jqP4VGy4ubSEOvB+XRCdTjWUJEuCkkWRyiRlxdKCNqjP8cTjRUg0JGhDYsW9RprSsQPqdnlOviWD11mpmynwSJzlRyzzT8rgCka5XEnLzq8="
      },
      "s3": {
        "s3SchemaVersion": "1.0",
        "configurationId": "SQS-Event",
        "bucket": {
          "name": "test",
          "ownerIdentity": {
            "principalId": "A2B5KBXGR14B9R"
          },
          "arn": "arn:aws:s3:::20220807-sqs-test"
        },
        "object": {
          "key": "eicar.txt",
          "size": 9846,
          "eTag": "ad1cdeed43375dca5b5e892be0968525",
          "sequencer": "0062EFCD57CFFC5419"
        }
      }
    }]}'
```
</details>

<details>
    <summary>test.txtファイルのアップロードイベントテスト用コマンド</summary>

```
aws sqs send-message --queue-url http://localhost:9324/000000000000/queue1 \
 --endpoint-url http://localhost:9324 \
 --message-body '{"Records": [{
      "eventVersion": "2.1",
      "eventSource": "aws:s3",
      "awsRegion": "ap-northeast-1",
      "eventTime": "2022-08-07T14:33:59.870Z",
      "eventName": "ObjectCreated:Put",
      "userIdentity": {
        "principalId": "AWS:AIDAVMRY2N7OKTN33RYNV"
      },
      "requestParameters": {
        "sourceIPAddress": "60.95.0.122"
      },
      "responseElements": {
        "x-amz-request-id": "Q73VJ1CPJ64CKJQ0",
        "x-amz-id-2": "jqP4VGy4ubSEOvB+XRCdTjWUJEuCkkWRyiRlxdKCNqjP8cTjRUg0JGhDYsW9RprSsQPqdnlOviWD11mpmynwSJzlRyzzT8rgCka5XEnLzq8="
      },
      "s3": {
        "s3SchemaVersion": "1.0",
        "configurationId": "SQS-Event",
        "bucket": {
          "name": "test",
          "ownerIdentity": {
            "principalId": "A2B5KBXGR14B9R"
          },
          "arn": "arn:aws:s3:::20220807-sqs-test"
        },
        "object": {
          "key": "test.txt",
          "size": 9846,
          "eTag": "ad1cdeed43375dca5b5e892be0968525",
          "sequencer": "0062EFCD57CFFC5419"
        }
      }
    }]}'
```

</details>


<details>
    <summary>キューのメッセージをすべてパージするコマンド</summary>

```
aws sqs purge-queue \
    --queue-url http://localhost:9324/000000000000/queue1 \
    --endpoint-url http://localhost:9324
```
</details>