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

clamd # clamdを起動