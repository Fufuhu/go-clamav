#!/bin/bash

# デフォルトの実行周期を3600秒に設定
DEFAULT_INTERVAL=3600

# 環境変数FRESHCLAM_INTERVALが設定されていない場合はデフォルト値を使用
INTERVAL=${FRESHCLAM_INTERVAL:-$DEFAULT_INTERVAL}

# Fargateなどでボリュームマウントする際に権限が書き変わるので、再度権限を変更
chmod 777 /var/lib/clamav

# 定期的にfreshclamコマンドを実行する関数
run_freshclam() {
  while true; do
    echo "Running freshclam..."
    freshclam
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
