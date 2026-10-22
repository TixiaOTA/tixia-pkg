# HTTP Client Package

A robust HTTP client utility package with built-in retry logic, exponential backoff, and timeout management for Go applications.

## Overview

The HTTP Client package enables you to:

- Make HTTP requests with automatic retry logic
- Configure exponential backoff for retries
- Set timeouts and deadlines for requests
- Handle GET, POST JSON, and POST form requests
- Manage request headers and query parameters
- Buffer response bodies for retry scenarios

## Installation

```bash
go get github.com/TixiaOTA/tixia-pkg/http/client
```

## Configuration

### Configuration

Create an HTTP client with configuration:

```go
import (
    "time"
    "github.com/TixiaOTA/tixia-pkg/http/client"
)

timeout := 30 * time.Second
deadline := 2 * time.Minute

httpClient := client.NewHTTPClient(client.HTTPConfig{
    BaseURL: "https://api.example.com",
    Timeout: &timeout,
    Deadline: &deadline,
    Retry: &client.RetryConfig{
        MaxRetries:        3,
        InitialBackoff:   100 * time.Millisecond,
        MaxBackoff:       5 * time.Second,
        BackoffMultiplier: 2.0,
    },
})
```

### Configuration Fields

#### HTTPConfig

- **BaseURL** (string, required): Base URL for all requests
- **Timeout** (\*time.Duration, optional): Timeout for individual request attempts. Defaults to 30 seconds
- **Deadline** (\*time.Duration, optional): Deadline for the entire operation (including retries). If set, retries stop when deadline is reached
- **Retry** (\*RetryConfig, optional): Retry configuration. Defaults to 3 retries with exponential backoff

#### RetryConfig

- **MaxRetries** (int, optional): Maximum number of retry attempts. Defaults to 3. Set to 0 to disable retries
- **InitialBackoff** (time.Duration, optional): Initial backoff duration. Defaults to 100ms
- **MaxBackoff** (time.Duration, optional): Maximum backoff duration. Defaults to 5s
- **BackoffMultiplier** (float64, optional): Multiplier for exponential backoff. Defaults to 2.0

## Usage

### GET Request

```go
import (
    "context"
    "github.com/TixiaOTA/tixia-pkg/http/client"
)

ctx := context.Background()

headers := map[string]string{
    "Authorization": "Bearer token",
}

resp, err := httpClient.Get(ctx, client.GetRequest{
    Path:    "/api/users",
    Headers: &headers,
})

if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()

// Read response body
body, err := io.ReadAll(resp.Body)
```

### POST JSON Request

```go
type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

user := User{
    Name:  "John Doe",
    Email: "john@example.com",
}

headers := map[string]string{
    "Authorization": "Bearer token",
}

resp, err := httpClient.PostJSON(ctx, client.PostJSONRequest{
    Path:    "/api/users",
    Body:    user,
    Headers: &headers,
})

if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()
```

### POST Form Request

```go
import "net/url"

formData := url.Values{}
formData.Set("username", "johndoe")
formData.Set("password", "secret")

resp, err := httpClient.PostForm(ctx, client.PostFormRequest{
    Path: "/api/login",
    Body: formData,
})

if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()
```

### Dynamic Configuration

You can modify client configuration at runtime:

```go
// Change base URL
httpClient.SetBaseURL("https://new-api.example.com")

// Change timeout
httpClient.SetTimeout(60 * time.Second)

// Disable retries
httpClient.SetMaxRetries(0)

// Configure retry behavior
httpClient.SetRetry(client.RetryConfig{
    MaxRetries:        5,
    InitialBackoff:   200 * time.Millisecond,
    MaxBackoff:       10 * time.Second,
    BackoffMultiplier: 1.5,
})
```

## API Reference

### Public Functions

#### NewHTTPClient

Creates a new HTTP client with the specified configuration.

```go
func NewHTTPClient(config HTTPConfig) *HTTPClient
```

**Parameters:**

- `config` (HTTPConfig, required): HTTP client configuration

**Return Values:**

- `*HTTPClient`: HTTP client instance

**Special Notes:** If RetryConfig is not provided, defaults to 3 retries with exponential backoff.

### Public Types

#### HTTPClient

HTTP client struct that manages requests and retries.

```go
type HTTPClient struct {
    config HTTPConfig
}
```

#### Methods on HTTPClient

##### SetBaseURL

Sets the base URL for the HTTP client.

```go
func (c *HTTPClient) SetBaseURL(baseURL string)
```

**Parameters:**

- `baseURL` (string, required): Base URL for all requests

##### SetTimeout

Sets the timeout for individual request attempts.

```go
func (c *HTTPClient) SetTimeout(timeout time.Duration)
```

**Parameters:**

- `timeout` (time.Duration, required): Timeout duration

##### SetDeadline

Sets the deadline for the entire operation (including retries).

```go
func (c *HTTPClient) SetDeadline(deadline time.Duration)
```

**Parameters:**

- `deadline` (time.Duration, required): Deadline duration

##### SetMaxRetries

Sets the maximum number of retries.

```go
func (c *HTTPClient) SetMaxRetries(maxRetries int)
```

**Parameters:**

- `maxRetries` (int, required): Maximum retry attempts. Set to 0 to disable retries

##### Get

Makes a GET request to the given path.

```go
func (c *HTTPClient) Get(ctx context.Context, request GetRequest) (*http.Response, error)
```

**Parameters:**

- `ctx` (context.Context, required): Context for the request
- `request` (GetRequest, required): GET request configuration

**Return Values:**

- `*http.Response`: HTTP response
- `error`: Error if request fails

**Special Notes:** Base URL is prepended to the path. Response body is buffered to allow retries.

##### PostJSON

Makes a POST request with JSON body to the given path.

```go
func (c *HTTPClient) PostJSON(ctx context.Context, request PostJSONRequest) (*http.Response, error)
```

**Parameters:**

- `ctx` (context.Context, required): Context for the request
- `request` (PostJSONRequest, required): POST JSON request configuration

**Return Values:**

- `*http.Response`: HTTP response
- `error`: Error if request fails

**Special Notes:** Body is automatically marshaled to JSON. Content-Type header is set to "application/json". Request body can be re-read for retries.

##### PostForm

Makes a POST request with form body to the given path.

```go
func (c *HTTPClient) PostForm(ctx context.Context, request PostFormRequest) (*http.Response, error)
```

**Parameters:**

- `ctx` (context.Context, required): Context for the request
- `request` (PostFormRequest, required): POST form request configuration

**Return Values:**

- `*http.Response`: HTTP response
- `error`: Error if request fails

**Special Notes:** Content-Type header is set to "application/x-www-form-urlencoded". Request body can be re-read for retries.

##### PostRaw

Makes a POST request with a raw body to the given path.

```go
func (c *HTTPClient) PostRaw(ctx context.Context, request PostRawRequest) (*http.Response, error)
```

**Parameters:**

- `ctx` (context.Context, required): Context for the request
- `request` (PostRawRequest, required): POST raw request configuration

**Return Values:**

- `*http.Response`: HTTP response
- `error`: Error if request fails

**Special Notes:** The caller is responsible for setting `ContentType` (for example `multipart/form-data`). Request body can be re-read for retries.

## Error Handling

The client handles errors gracefully:

- **Network errors**: Retries on timeout and temporary connection errors
- **HTTP 5xx errors**: Automatically retries on server errors (500-599)
- **Deadline exceeded**: Stops retrying when deadline is reached
- **Max retries exceeded**: Returns the last error/response after all retries are exhausted

### Retry Behavior

The client retries requests in the following scenarios:

- Network timeout errors
- Temporary connection errors
- HTTP 5xx server errors
- Context deadline exceeded (if deadline is set)

Retries use exponential backoff:

- First retry: InitialBackoff duration
- Subsequent retries: Previous backoff × BackoffMultiplier
- Maximum backoff capped at MaxBackoff

## Data Models

### HTTPConfig

HTTP client configuration:

```go
type HTTPConfig struct {
    BaseURL  string
    Timeout  *time.Duration
    Deadline *time.Duration
    Retry    *RetryConfig
}
```

### RetryConfig

Retry configuration:

```go
type RetryConfig struct {
    MaxRetries        int
    InitialBackoff    time.Duration
    MaxBackoff        time.Duration
    BackoffMultiplier float64
}
```

### GetRequest

GET request configuration:

```go
type GetRequest struct {
    Path    string
    Headers *map[string]string
}
```

### PostJSONRequest

POST JSON request configuration:

```go
type PostJSONRequest struct {
    Path    string
    Body    any
    Headers *map[string]string
}
```

### PostFormRequest

POST form request configuration:

```go
type PostFormRequest struct {
    Path    string
    Body    url.Values
    Headers *map[string]string
}
```

### PostRawRequest

POST raw request configuration:

```go
type PostRawRequest struct {
    Path        string
    Body        []byte
    ContentType string
    Headers     *map[string]string
}
```

## Examples

### Complete Example

```go
package main

import (
    "context"
    "fmt"
    "io"
    "log"
    "time"

    "github.com/TixiaOTA/tixia-pkg/http/client"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func main() {
    // Create HTTP client
    timeout := 30 * time.Second
    deadline := 2 * time.Minute

    httpClient := client.NewHTTPClient(client.HTTPConfig{
        BaseURL: "https://api.example.com",
        Timeout: &timeout,
        Deadline: &deadline,
        Retry: &client.RetryConfig{
            MaxRetries:        3,
            InitialBackoff:   100 * time.Millisecond,
            MaxBackoff:       5 * time.Second,
            BackoffMultiplier: 2.0,
        },
    })

    ctx := context.Background()

    // GET request
    headers := map[string]string{
        "Authorization": "Bearer token123",
    }

    resp, err := httpClient.Get(ctx, client.GetRequest{
        Path:    "/api/users/1",
        Headers: &headers,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Response: %s\n", string(body))

    // POST JSON request
    newUser := User{
        Name:  "Jane Doe",
        Email: "jane@example.com",
    }

    resp, err = httpClient.PostJSON(ctx, client.PostJSONRequest{
        Path:    "/api/users",
        Body:    newUser,
        Headers: &headers,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    fmt.Printf("Created user, status: %d\n", resp.StatusCode)
}
```

### Example with Dynamic Configuration

```go
// Create client with defaults
httpClient := client.NewHTTPClient(client.HTTPConfig{
    BaseURL: "https://api.example.com",
})

// Modify configuration at runtime
httpClient.SetTimeout(60 * time.Second)
httpClient.SetMaxRetries(5)

// Make request
resp, err := httpClient.Get(ctx, client.GetRequest{
    Path: "/api/data",
})
```

## License

This package is part of the TixiaOTA project.
