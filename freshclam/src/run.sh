#!/bin/bash

# デフォルトの実行周期を3600秒に設定
DEFAULT_INTERVAL=3600

# 環境変数FRESHCLAM_INTERVALが設定されていない場合はデフォルト値を使用
INTERVAL=${FRESHCLAM_INTERVAL:-$DEFAULT_INTERVAL}

# read-only root filesystem 対応:
# 設定ファイルの生成先である /run/clamav を用意する。
# Fargate 等でボリュームをマウントすると所有者が root に戻るため、起動のたびに chown する。
mkdir -p /run/clamav
chown clamav:clamav /run/clamav

# 設定ファイルは read-only な /etc/clamav ではなく、書き込み可能な /run/clamav に生成する
envsubst < /etc/clamav/freshclam.conf.template > /run/clamav/freshclam.conf

# Fargateなどでボリュームマウントする際に権限が書き変わるので、再度権限を変更
chmod 777 /var/lib/clamav

# 定期的にfreshclamコマンドを実行する関数
run_freshclam() {
  while true; do
    echo "Running freshclam..."
    freshclam --config-file=/run/clamav/freshclam.conf
    if [ $? -ne 0 ]; then
      echo "freshclam command failed. Exiting script."
      exit 1
    fi
    echo "Waiting for $INTERVAL seconds before the next run..."
    sleep $INTERVAL
  done
}

# 実行
run_freshclam
