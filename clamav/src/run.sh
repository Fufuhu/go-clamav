#!/bin/bash

# read-only root filesystem 対応:
# 設定ファイルの生成先・ソケットの配置先である /run/clamav を用意する。
# Fargate 等でボリュームをマウントすると所有者が root に戻るため、起動のたびに chown する。
mkdir -p /run/clamav
chown clamav:clamav /run/clamav

# 設定ファイルは read-only な /etc/clamav ではなく、書き込み可能な /run/clamav に生成する
envsubst < /etc/clamav/clamd.conf.template > /run/clamav/clamd.conf

# freshclamがデータベースをダウンロードするまで待機
retry_count=0
while [ ! -f /var/lib/clamav/main.* ] || [ ! -f /var/lib/clamav/daily.* ] || [ ! -f /var/lib/clamav/bytecode.* ]; do
  echo "Waiting for freshclam to download all database files..."
  sleep 10
  retry_count=$((retry_count + 1))
  if [ $retry_count -gt 30 ]; then
    echo "Failed to download database files. Pleas check freshclam logs. Exiting script."
    exit 1
  fi
done

echo "All database files are ready. Starting clamd..."

clamd -c /run/clamav/clamd.conf # clamdを起動
