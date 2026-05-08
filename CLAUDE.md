# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 概要

go-clamav は、ClamAV を用いて S3 にアップロードされたファイルをウイルススキャンする Go 製サービスです。S3 イベント通知が送られた SQS キューをロングポーリングし、対象オブジェクトを取得して TCP の `INSTREAM` プロトコルで clamd にストリーミング送信し、結果を DynamoDB に保存します。ローカル開発環境は S3 / SQS / DynamoDB 互換の代替ミドルウェアでフルにコンテナ化されています。

## 主なコマンド

Go のバージョンは `.go-version`（現在 1.22.3）に固定されています。ビルド前に goenv / asdf 等で揃えてください。

```bash
# バイナリをビルド（出力: ./go-clamav）
go build

# ローカルインフラを起動: clamd + freshclam + minio + elasticmq + dynamodb (+ 各管理 UI)
# 初回起動時に DynamoDB のテーブル作成が間に合わないことがある
# （aws-cli サイドカーが dynamodb コンテナとレースする）。その場合は再実行する。
docker compose up -d

# docker compose スタックに対してポーラーをローカル実行
./exec_local.sh        # ローカル用の環境変数を付けて `./go-clamav poll` を実行する

# 全テスト実行
go test ./...

# 単一パッケージのテスト実行
go test ./internal/virus_scan/scanner/...

# 単一テスト名で実行
go test ./config -run TestGetConfig
```

各種ローカル UI: minio `http://localhost:9001`、ElasticMQ `http://localhost:9325`、DynamoDB Admin `http://localhost:8001`。

スキャンパイプラインの動作確認は、MinIO の `test` バケットに `eicar.txt` と `test.txt` をアップロードし、ElasticMQ に S3 イベント形式の SQS メッセージを送り込むことで行います。具体的な `aws sqs send-message` のペイロードは README.md を参照してください（MinIO → ElasticMQ のイベント連携は構成されていないため、SQS は手動で叩く必要があります）。

## アーキテクチャ

### 実行フロー（コマンドは `go-clamav poll` のみ）
1. `main.go` → `cmd.Execute()`（cobra）→ `pollCmd` → `internal/cmd/poll.CommandPoll.Run`。
2. `Run` は `config.Configuration` から 4 つのクライアント（SQS / S3 / DynamoDB / clamd）を生成し、`scanner.Scanner` に注入する。
3. `sqs.Client.Poll` は無限ループで動作し、`ReceiveMessages` → 各メッセージに対して `process(...)`（= `scanner.Process`）→ 成功時に `DeleteMessage` を実行する。エラー時はログのみで、メッセージは再配信に任せる（アプリ側に DLQ の仕掛けはない）。
4. `scanner.Process` の流れ: `SCANNING_TARGET_FILE_PATTERNS` 正規表現でフィルタ（`message.IsTargetFile`）→ `s3.GetObject` → `clamav.Scan`（チャンク単位で clamd にストリーム）→ `dynamodb.PutScanResult`。感染ファイルの場合は `DYNAMODB_TABLE` と `DYNAMODB_TABLE_INFECTED` の **両方** に書き込まれる。

### レイヤ構成
- `cmd/` — cobra のエントリポイント。薄く保ち、`internal/cmd/<subcommand>` に処理を委譲する。新しいサブコマンドを足すときは `internal/cmd.CommandInterface` を実装し、`cmd/` で登録する。
- `config/` — `kelseyhightower/envconfig` で環境変数から生成される単一の `Configuration` 構造体。`GetConfig()` はメモ化される。テストで使う `Initialize()` でキャッシュをクリアできる。
- `internal/queue/clients/` — スキャナが依存する抽象は `QueueMessageInterface`。具象構造体 `QueueMessage` は複数パッケージで共有されており、`model.ScanResult` にも埋め込まれている。SQS の `Client.Poll` 内で S3 イベント JSON をパースし、`Records[]` 1 件ごとに `QueueMessage` を 1 つ生成する。
- `internal/virus_scan/clients/clamav/` — TCP `INSTREAM` プロトコルを直接実装している（サードパーティの clamav ライブラリは使用していない）。チャンクサイズは 1024 バイト。結果は clamd の生レスポンス文字列で、`ResultOK`（`"stream: OK\n"`）と比較して判定する。
- `internal/objects/clients/s3/`、`internal/db/clients/dynamodb/` — AWS SDK v2 の薄いラッパ。SQS / S3 / DynamoDB の各クライアントはいずれも `*BaseUrl` 設定でローカルエミュレータ（MinIO / ElasticMQ / DynamoDB Local）にエンドポイントを差し替えできる。
- `internal/model/ScanResult` は `clients.QueueMessage` を埋め込んでおり、DynamoDB への書き込みもキュー処理と同じ Bucket / Key / ObjectPath メソッドを使う。DynamoDB のハッシュキーは `ObjectPath`（`s3://bucket/key`）なので、同一オブジェクトを再スキャンすると **既存行が上書きされる** 点に注意。

### サイドカー / インフラ要素（Go コード以外）
- `clamav/` と `freshclam/` — clamd デーモンとそのシグネチャ更新コンテナの Dockerfile。`compose.yaml` 上で `signatures` という名前付きボリュームを共有している。
- `signature_downloader/` — ClamAV の公式シグネチャをプライベート S3 バケットへミラーするための、ECS 等での定期実行を想定した独立したツール。`freshclam` を `database.clamav.net` ではなくプライベートミラーから引かせる構成で利用する。独自の `docker-compose.yaml` と README を持つ。

### CI ワークフロー（`.github/workflows/`）
- `build_and_push.yml` — リポジトリ直下の `Dockerfile`（go-clamav 本体）をビルドし、`fufuhu/go-clamav:<ref>-<short-sha>` として Docker Hub に push する。`workflow_dispatch` 起動で、`build_target` 入力でブランチ / タグを指定できる。
- `build_and_push_signiture_downloader.yml` — `signature_downloader/` のイメージをビルドし、`fufuhu/go-clamav-signiture-downloader:<ref>-<short-sha>` として push する（`context: signature_downloader` を指定。タグ名は履歴上の表記揺れにより `signiture` のまま）。
