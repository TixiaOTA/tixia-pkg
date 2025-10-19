package client

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"go.uber.org/zap"
)

type S3Client struct {
	client     *awss3.Client
	config     S3Config
	logger     *zap.Logger
	bucketName string
}

func NewS3Client(config S3Config, logger *zap.Logger) (*S3Client, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	awsConfig, err := config.loadAWSConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client with custom endpoint if needed
	var s3Client *awss3.Client
	if config.Endpoint != AwsEndpoint {
		s3Client = awss3.NewFromConfig(awsConfig, func(o *awss3.Options) {
			o.BaseEndpoint = aws.String(config.Endpoint)
		})
	} else {
		s3Client = awss3.NewFromConfig(awsConfig)
	}

	client := &S3Client{
		client:     s3Client,
		config:     config,
		logger:     logger,
		bucketName: config.BucketName,
	}

	return client, nil
}

func (c *S3Config) loadAWSConfig() (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			c.AccessKeyID,
			c.SecretAccessKey,
			"",
		)),
		config.WithRegion(c.Region),
	)
	if err != nil {
		return aws.Config{}, err
	}

	return cfg, nil
}

func (c *S3Client) UploadFile(ctx context.Context, key string, file io.Reader, contentType string) error {
	if key == "" {
		return ErrInvalidKey
	}

	c.logger.Info("uploading file", zap.String("key", key), zap.String("contentType", contentType))

	_, err := c.client.PutObject(ctx, &awss3.PutObjectInput{
		Bucket:      aws.String(c.bucketName),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(contentType),
	})

	if err != nil {
		c.logger.Error("failed to upload file", zap.String("key", key), zap.Error(err))
		return fmt.Errorf("%w: %v", ErrUploadFailed, err)
	}

	c.logger.Info("file uploaded successfully", zap.String("key", key))
	return nil
}

func (c *S3Client) UploadFilePublic(ctx context.Context, key string, file io.Reader, contentType string) error {
	if key == "" {
		return ErrInvalidKey
	}

	c.logger.Info("uploading file", zap.String("key", key), zap.String("contentType", contentType))

	_, err := c.client.PutObject(ctx, &awss3.PutObjectInput{
		Bucket:      aws.String(c.bucketName),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(contentType),
		ACL:         types.ObjectCannedACLPublicRead,
	})

	if err != nil {
		c.logger.Error("failed to upload file", zap.String("key", key), zap.Error(err))
		return fmt.Errorf("%w: %v", ErrUploadFailed, err)
	}

	c.logger.Info("file uploaded successfully", zap.String("key", key))
	return nil
}

func (c *S3Client) GetPresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	if key == "" {
		return "", ErrInvalidKey
	}

	if expiry <= 0 {
		expiry = c.config.PresignExpiry
	}

	c.logger.Info("generating presigned URL", zap.String("key", key), zap.Duration("expiry", expiry))

	presignClient := awss3.NewPresignClient(c.client)
	request, err := presignClient.PresignGetObject(ctx, &awss3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	}, func(opts *awss3.PresignOptions) {
		opts.Expires = expiry
	})

	if err != nil {
		c.logger.Error("failed to generate presigned URL", zap.String("key", key), zap.Error(err))
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return request.URL, nil
}

func (c *S3Client) GetPublicURL(ctx context.Context, key string) (string, error) {
	if key == "" {
		return "", ErrInvalidKey
	}

	c.logger.Info("getting public URL", zap.String("key", key))

	return fmt.Sprintf("https://%s.%s/%s", c.config.BucketName, c.config.BaseURL, key), nil
}

func (c *S3Client) DeleteFile(ctx context.Context, key string) error {
	if key == "" {
		return ErrInvalidKey
	}

	c.logger.Info("deleting file", zap.String("key", key))

	_, err := c.client.DeleteObject(ctx, &awss3.DeleteObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		c.logger.Error("failed to delete file", zap.String("key", key), zap.Error(err))
		return fmt.Errorf("%w: %v", ErrDeleteFailed, err)
	}

	c.logger.Info("file deleted successfully", zap.String("key", key))
	return nil
}
