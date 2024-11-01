#!/bin/bash

echo "Start ClamAV signatures uploaded to S3."

# ClamAVデータベースのダウンロード
freshclam --datadir=/tmp/clamdb

# S3へのアップロード
aws s3 cp /tmp/clamdb/main.cvd s3://$MIRROR_SITE_S3_BUCKET_NAME/main.cvd --storage-class STANDARD_IA
aws s3 cp /tmp/clamdb/daily.cvd s3://$MIRROR_SITE_S3_BUCKET_NAME/daily.cvd --storage-class STANDARD_IA
aws s3 cp /tmp/clamdb/bytecode.cvd s3://$MIRROR_SITE_S3_BUCKET_NAME/bytecode.cvd --storage-class STANDARD_IA

echo "ClamAV signatures uploaded to S3 end."
