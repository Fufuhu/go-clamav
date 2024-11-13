QUEUE_URL=http://localhost:9324/000000000000/queue1 \
BASE_URL=http://localhost:9324 \
DYNAMODB_BASE_URL=http://localhost:8000 \
S3_BASE_URL=http://localhost:9000 \
AWS_ACCESS_KEY_ID=minio \
AWS_SECRET_ACCESS_KEY=miniominio \
SCANNING_TARGET_FILE_PATTERNS=".*-raw\..*,test.txt" \
./go-clamav poll