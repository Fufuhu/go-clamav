package clients

import "fmt"

type S3Object struct {
	Bucket string
	Key    string
}

func (s3Object *S3Object) GetObjectPath() string {
	path := fmt.Sprintf("s3://%s/%s", s3Object.Bucket, s3Object.Key)
	return path
}
