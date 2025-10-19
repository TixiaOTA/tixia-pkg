package client

import (
	"errors"
	"time"
)

var (
	ErrMissingEndpoint        = errors.New("missing endpoint")
	ErrMissingRegion          = errors.New("missing region")
	ErrMissingAccessKeyID     = errors.New("missing access key ID")
	ErrMissingSecretAccessKey = errors.New("missing secret access key")
	ErrMissingBucketName      = errors.New("missing bucket name")
	ErrMissingPresignExpiry   = errors.New("missing presign expiry")
	ErrUploadFailed           = errors.New("upload failed")
	ErrDeleteFailed           = errors.New("delete failed")
	ErrInvalidKey             = errors.New("invalid key")
	ErrMissingBaseURL         = errors.New("missing base url")
	AwsEndpoint               = "s3.amazonaws.com"
)

type S3Config struct {
	Endpoint        string        `json:"endpoint"`
	Region          string        `json:"region"`
	AccessKeyID     string        `json:"access_key_id"`
	SecretAccessKey string        `json:"secret_access_key"`
	BucketName      string        `json:"bucket_name"`
	UseSSL          bool          `json:"use_ssl" default:"true"`
	PresignExpiry   time.Duration `json:"presign_expiry"`
	BaseURL         string        `json:"base_url"`
}

func (c *S3Config) Validate() error {
	if c.Endpoint == "" {
		return ErrMissingEndpoint
	}
	if c.Region == "" {
		return ErrMissingRegion
	}
	if c.AccessKeyID == "" {
		return ErrMissingAccessKeyID
	}
	if c.SecretAccessKey == "" {
		return ErrMissingSecretAccessKey
	}
	if c.BucketName == "" {
		return ErrMissingBucketName
	}
	if c.PresignExpiry == 0 {
		return ErrMissingPresignExpiry
	}
	if c.BaseURL == "" {
		return ErrMissingBaseURL
	}
	return nil
}
