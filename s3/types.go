package s3

import (
	"context"
	"io"
	"time"
)

type S3Client interface {
	UploadFile(ctx context.Context, key string, file io.Reader, contentType string) error
	GetPresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error)
	DeleteFile(ctx context.Context, key string) error
}
