# S3 Package

A Go client library for AWS S3 and S3-compatible object storage services, providing file upload, download, and URL generation capabilities.

## Overview

The S3 package enables you to:
- Upload files to S3 buckets
- Upload files with public read access
- Generate presigned URLs for secure file access
- Generate public URLs for publicly accessible files
- Delete files from S3 buckets
- Support custom S3-compatible endpoints (e.g., MinIO, DigitalOcean Spaces)
- Use AWS SDK v2 for modern AWS integration

## Installation

```bash
go get github.com/TixiaOTA/tixia-pkg/s3
```

## Configuration

### Configuration

Create an S3 client with configuration:

```go
import (
    "time"
    "github.com/TixiaOTA/tixia-pkg/s3/client"
    "go.uber.org/zap"
)

config := client.S3Config{
    Endpoint:        "s3.amazonaws.com", // or custom endpoint
    Region:          "us-east-1",
    AccessKeyID:     "your-access-key-id",
    SecretAccessKey: "your-secret-access-key",
    BucketName:      "my-bucket",
    UseSSL:          true,
    PresignExpiry:   1 * time.Hour,
    BaseURL:         "s3.amazonaws.com", // for public URLs
}

logger, _ := zap.NewProduction()

s3Client, err := client.NewS3Client(config, logger)
if err != nil {
    log.Fatal(err)
}
```

### Configuration Fields

#### S3Config

- **Endpoint** (string, required): S3 endpoint URL. Use `"s3.amazonaws.com"` for AWS S3, or custom endpoint for S3-compatible services
- **Region** (string, required): AWS region (e.g., `"us-east-1"`, `"eu-west-1"`)
- **AccessKeyID** (string, required): AWS access key ID
- **SecretAccessKey** (string, required): AWS secret access key
- **BucketName** (string, required): S3 bucket name
- **UseSSL** (bool, optional): Enable SSL/TLS. Defaults to true
- **PresignExpiry** (time.Duration, required): Default expiry time for presigned URLs
- **BaseURL** (string, required): Base URL for generating public URLs (e.g., `"s3.amazonaws.com"` or `"nyc3.digitaloceanspaces.com"`)

## Usage

### Upload File

```go
import (
    "context"
    "os"
    "github.com/TixiaOTA/tixia-pkg/s3/client"
)

ctx := context.Background()

// Open file
file, err := os.Open("path/to/file.jpg")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

// Upload file
err = s3Client.UploadFile(ctx, "uploads/file.jpg", file, "image/jpeg")
if err != nil {
    log.Fatal(err)
}
```

### Upload File with Public Read Access

```go
// Upload file with public read ACL
err = s3Client.UploadFilePublic(ctx, "public/file.jpg", file, "image/jpeg")
if err != nil {
    log.Fatal(err)
}
```

### Generate Presigned URL

```go
import "time"

// Generate presigned URL with default expiry
url, err := s3Client.GetPresignedURL(ctx, "uploads/file.jpg", 0)
if err != nil {
    log.Fatal(err)
}

// Generate presigned URL with custom expiry
url, err := s3Client.GetPresignedURL(ctx, "uploads/file.jpg", 24*time.Hour)
if err != nil {
    log.Fatal(err)
}

// Use URL to access file securely
fmt.Printf("Presigned URL: %s\n", url)
```

### Generate Public URL

```go
// Generate public URL (for files uploaded with public read access)
url, err := s3Client.GetPublicURL(ctx, "public/file.jpg")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Public URL: %s\n", url)
```

### Delete File

```go
err = s3Client.DeleteFile(ctx, "uploads/file.jpg")
if err != nil {
    log.Fatal(err)
}
```

### Using S3-Compatible Services

```go
// Example: DigitalOcean Spaces
config := client.S3Config{
    Endpoint:        "nyc3.digitaloceanspaces.com",
    Region:          "nyc3",
    AccessKeyID:     "your-spaces-key",
    SecretAccessKey: "your-spaces-secret",
    BucketName:      "my-space",
    UseSSL:          true,
    PresignExpiry:   1 * time.Hour,
    BaseURL:         "nyc3.digitaloceanspaces.com",
}

s3Client, err := client.NewS3Client(config, logger)
```

## API Reference

### Public Functions

#### NewS3Client

Creates a new S3 client with the specified configuration.

```go
func NewS3Client(config S3Config, logger *zap.Logger) (*S3Client, error)
```

**Parameters:**

- `config` (S3Config, required): S3 configuration
- `logger` (*zap.Logger, required): Logger instance

**Return Values:**

- `*S3Client`: S3 client instance
- `error`: Error if configuration is invalid or AWS config fails

**Special Notes:** Validates configuration before creating client. Supports custom endpoints for S3-compatible services.

### Public Types

#### S3Client

S3 client struct for object storage operations.

```go
type S3Client struct {
    client     *awss3.Client
    config     S3Config
    logger     *zap.Logger
    bucketName string
}
```

#### Methods on S3Client

##### UploadFile

Uploads a file to S3 bucket.

```go
func (c *S3Client) UploadFile(ctx context.Context, key string, file io.Reader, contentType string) error
```

**Parameters:**

- `ctx` (context.Context, required): Context
- `key` (string, required): Object key (file path in bucket)
- `file` (io.Reader, required): File content reader
- `contentType` (string, required): Content type (MIME type)

**Return Values:**

- `error`: Error if upload fails

**Special Notes:** Key cannot be empty. Logs upload progress.

##### UploadFilePublic

Uploads a file to S3 bucket with public read access.

```go
func (c *S3Client) UploadFilePublic(ctx context.Context, key string, file io.Reader, contentType string) error
```

**Parameters:**

- `ctx` (context.Context, required): Context
- `key` (string, required): Object key
- `file` (io.Reader, required): File content reader
- `contentType` (string, required): Content type

**Return Values:**

- `error`: Error if upload fails

**Special Notes:** Sets ACL to public-read. File will be publicly accessible.

##### GetPresignedURL

Generates a presigned URL for secure file access.

```go
func (c *S3Client) GetPresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error)
```

**Parameters:**

- `ctx` (context.Context, required): Context
- `key` (string, required): Object key
- `expiry` (time.Duration, optional): URL expiry time. If <= 0, uses PresignExpiry from config

**Return Values:**

- `string`: Presigned URL
- `error`: Error if key is invalid or URL generation fails

**Special Notes:** Presigned URLs allow temporary access without exposing credentials.

##### GetPublicURL

Generates a public URL for publicly accessible files.

```go
func (c *S3Client) GetPublicURL(ctx context.Context, key string) (string, error)
```

**Parameters:**

- `ctx` (context.Context, required): Context
- `key` (string, required): Object key

**Return Values:**

- `string`: Public URL
- `error`: Error if key is invalid

**Special Notes:** Only works for files uploaded with public read access.

##### DeleteFile

Deletes a file from S3 bucket.

```go
func (c *S3Client) DeleteFile(ctx context.Context, key string) error
```

**Parameters:**

- `ctx` (context.Context, required): Context
- `key` (string, required): Object key to delete

**Return Values:**

- `error`: Error if deletion fails

**Special Notes:** Key cannot be empty. Logs deletion progress.

## Error Handling

The package defines specific error types:

- **ErrMissingEndpoint**: Endpoint not provided in config
- **ErrMissingRegion**: Region not provided in config
- **ErrMissingAccessKeyID**: Access key ID not provided
- **ErrMissingSecretAccessKey**: Secret access key not provided
- **ErrMissingBucketName**: Bucket name not provided
- **ErrMissingPresignExpiry**: Presign expiry not provided
- **ErrMissingBaseURL**: Base URL not provided
- **ErrInvalidKey**: Invalid or empty key provided
- **ErrUploadFailed**: File upload failed
- **ErrDeleteFailed**: File deletion failed

**Example:**

```go
err := s3Client.UploadFile(ctx, key, file, contentType)
if err != nil {
    if errors.Is(err, client.ErrInvalidKey) {
        log.Println("Invalid key provided")
    } else if errors.Is(err, client.ErrUploadFailed) {
        log.Println("Upload failed, check AWS credentials and permissions")
    } else {
        log.Printf("Unexpected error: %v", err)
    }
}
```

## Data Models

### S3Config

S3 client configuration:

```go
type S3Config struct {
    Endpoint        string
    Region          string
    AccessKeyID     string
    SecretAccessKey string
    BucketName      string
    UseSSL          bool
    PresignExpiry   time.Duration
    BaseURL         string
}
```

### S3Client Interface

The package also provides an interface for S3 operations:

```go
type S3Client interface {
    UploadFile(ctx context.Context, key string, file io.Reader, contentType string) error
    GetPresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error)
    DeleteFile(ctx context.Context, key string) error
}
```

## Examples

### Complete Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/TixiaOTA/tixia-pkg/s3/client"
    "go.uber.org/zap"
)

func main() {
    // Initialize logger
    logger, _ := zap.NewProduction()

    // Configure S3 client
    config := client.S3Config{
        Endpoint:        "s3.amazonaws.com",
        Region:          "us-east-1",
        AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
        SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
        BucketName:      "my-bucket",
        UseSSL:          true,
        PresignExpiry:   1 * time.Hour,
        BaseURL:         "s3.amazonaws.com",
    }

    // Create S3 client
    s3Client, err := client.NewS3Client(config, logger)
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Upload file
    file, err := os.Open("image.jpg")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    err = s3Client.UploadFile(ctx, "uploads/image.jpg", file, "image/jpeg")
    if err != nil {
        log.Fatal(err)
    }

    // Generate presigned URL
    url, err := s3Client.GetPresignedURL(ctx, "uploads/image.jpg", 24*time.Hour)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("File uploaded. Access URL: %s\n", url)

    // Delete file
    err = s3Client.DeleteFile(ctx, "uploads/image.jpg")
    if err != nil {
        log.Fatal(err)
    }
}
```

### Example with Public Upload

```go
// Upload file with public access
file, _ := os.Open("public-image.jpg")
defer file.Close()

err = s3Client.UploadFilePublic(ctx, "public/image.jpg", file, "image/jpeg")
if err != nil {
    log.Fatal(err)
}

// Generate public URL
publicURL, err := s3Client.GetPublicURL(ctx, "public/image.jpg")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Public URL: %s\n", publicURL)
```

### Example with MinIO

```go
// Configure for MinIO
config := client.S3Config{
    Endpoint:        "localhost:9000",
    Region:          "us-east-1",
    AccessKeyID:     "minioadmin",
    SecretAccessKey: "minioadmin",
    BucketName:      "my-bucket",
    UseSSL:          false, // MinIO may not use SSL locally
    PresignExpiry:   1 * time.Hour,
    BaseURL:         "localhost:9000",
}

s3Client, err := client.NewS3Client(config, logger)
```

## License

This package is part of the TixiaOTA project.

