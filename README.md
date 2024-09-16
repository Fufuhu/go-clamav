

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

```
aws sqs purge-queue \
    --queue-url http://localhost:9324/000000000000/queue1 \
    --endpoint-url http://localhost:9324
```