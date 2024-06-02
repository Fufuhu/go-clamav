package clients

import "context"

type ClientInterface interface {
	Poll(ctx context.Context) ([]S3Object, error)
}
