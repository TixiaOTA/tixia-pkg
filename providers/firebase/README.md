# Firebase Package

A Go client library for Firebase Authentication and Firebase Cloud Messaging (FCM) services. This package provides methods for user management and sending push notifications to individual devices.

## Overview

The Firebase package enables you to:

- Manage Firebase Authentication users (create, get, update, delete)
- Verify ID tokens and create custom tokens
- Send push notifications to individual devices via FCM
- Handle platform-specific notification configurations (Android, iOS, Web)

## Installation

```bash
go get github.com/TixiaOTA/tixia-pkg/providers/firebase
```

## Configuration

### Client Configuration

Create a client instance by providing your Firebase project ID and service account JSON:

```go
import (
    "os"
    "github.com/TixiaOTA/tixia-pkg/providers/firebase"
    "github.com/TixiaOTA/tixia-pkg/logger"
)

// Read service account JSON from file or environment
serviceAccountJSON, err := os.ReadFile("path/to/service-account.json")
if err != nil {
    log.Fatal(err)
}

// Create logger
logger := logger.New()

// Create Firebase client
config := firebase.Config{
    ProjectID:          "your-project-id",
    ServiceAccountJSON: serviceAccountJSON,
    Logger:             logger,
}

client, err := firebase.New(config)
if err != nil {
    log.Fatal(err)
}
```

### Configuration Fields

- **ProjectID** (string, required): Firebase project ID
- **ServiceAccountJSON** ([]byte, required): Service account key JSON as byte array
- **Logger** (\*logger.Logger, required): tixia-pkg logger instance for request logging

**Note:** The package handles Firebase app initialization and authentication using the provided service account JSON. The service account must have appropriate permissions for Authentication and Cloud Messaging APIs.

## Usage

### Authentication

#### Create User

Create a new user account with email and password:

```go
ctx := context.Background()

req := firebase.CreateUserRequest{
    Email:       "user@example.com",
    Password:    "securepassword123",
    DisplayName: "John Doe",
    PhoneNumber: "+1234567890",
    Disabled:    false,
}

user, firebaseErr := client.CreateUser(ctx, req)
if firebaseErr != nil {
    log.Printf("Error creating user: %v", firebaseErr)
    return
}

fmt.Printf("Created user: %s (%s)\n", user.DisplayName, user.UID)
```

#### Get User

Retrieve user information by UID:

```go
req := firebase.GetUserRequest{
    UID: "user-uid-here",
}

user, firebaseErr := client.GetUser(ctx, req)
if firebaseErr != nil {
    log.Printf("Error getting user: %v", firebaseErr)
    return
}

fmt.Printf("User: %s, Email: %s, Verified: %t\n",
    user.DisplayName, user.Email, user.EmailVerified)
```

#### Get User by Email

Retrieve user information by email address:

```go
req := firebase.GetUserByEmailRequest{
    Email: "user@example.com",
}

user, firebaseErr := client.GetUserByEmail(ctx, req)
if firebaseErr != nil {
    log.Printf("Error getting user by email: %v", firebaseErr)
    return
}

fmt.Printf("Found user: %s\n", user.UID)
```

#### Update User

Update user profile information:

```go
req := firebase.UpdateUserRequest{
    UID:         "user-uid-here",
    DisplayName: stringPtr("Updated Name"),
    PhoneNumber: stringPtr("+1987654321"),
    Disabled:    boolPtr(false),
    CustomClaims: map[string]interface{}{
        "role": "admin",
        "premium": true,
    },
}

user, firebaseErr := client.UpdateUser(ctx, req)
if firebaseErr != nil {
    log.Printf("Error updating user: %v", firebaseErr)
    return
}

fmt.Printf("Updated user: %s\n", user.DisplayName)
```

#### Delete User

Delete a user account:

```go
req := firebase.DeleteUserRequest{
    UID: "user-uid-here",
}

firebaseErr := client.DeleteUser(ctx, req)
if firebaseErr != nil {
    log.Printf("Error deleting user: %v", firebaseErr)
    return
}

fmt.Println("User deleted successfully")
```

#### Verify ID Token

Verify and decode a Firebase ID token:

```go
req := firebase.VerifyIDTokenRequest{
    IDToken: "id-token-from-client",
}

token, firebaseErr := client.VerifyIDToken(ctx, req)
if firebaseErr != nil {
    log.Printf("Error verifying token: %v", firebaseErr)
    return
}

fmt.Printf("Token subject: %s, expires: %d\n", token.Subject, token.ExpiresAt)
```

#### Create Custom Token

Create a custom token for client-side authentication:

```go
req := firebase.CreateCustomTokenRequest{
    UID: "user-uid-here",
    CustomClaims: map[string]interface{}{
        "role": "user",
        "premium": false,
    },
}

customToken, firebaseErr := client.CreateCustomToken(ctx, req)
if firebaseErr != nil {
    log.Printf("Error creating custom token: %v", firebaseErr)
    return
}

fmt.Printf("Custom token: %s\n", customToken)
```

#### Revoke Refresh Tokens

Revoke all refresh tokens for a user (forces re-authentication):

```go
req := firebase.RevokeRefreshTokensRequest{
    UID: "user-uid-here",
}

firebaseErr := client.RevokeRefreshTokens(ctx, req)
if firebaseErr != nil {
    log.Printf("Error revoking tokens: %v", firebaseErr)
    return
}

fmt.Println("Refresh tokens revoked successfully")
```

### Firebase Cloud Messaging (FCM)

#### Send Notification to Device

Send a push notification to a specific device:

```go
req := firebase.SendToDeviceRequest{
    Message: firebase.Message{
        Token: "device-fcm-token-here",
        Notification: &firebase.Notification{
            Title: "Hello!",
            Body:  "This is a test notification",
            Image: "https://example.com/image.png",
        },
        Data: map[string]string{
            "click_action": "OPEN_ACTIVITY",
            "user_id": "12345",
        },
    },
}

response, firebaseErr := client.SendToDevice(ctx, req)
if firebaseErr != nil {
    log.Printf("Error sending message: %v", firebaseErr)
    return
}

fmt.Printf("Message sent successfully, ID: %s\n", response.Name)
```

#### Send Data Message

Send a data-only message (no notification):

```go
req := firebase.SendToDeviceRequest{
    Message: firebase.Message{
        Token: "device-fcm-token-here",
        Data: map[string]string{
            "type": "sync_request",
            "user_id": "12345",
            "action": "refresh_data",
        },
    },
}

response, firebaseErr := client.SendToDevice(ctx, req)
if firebaseErr != nil {
    log.Printf("Error sending data message: %v", firebaseErr)
    return
}

fmt.Printf("Data message sent, ID: %s\n", response.Name)
```

#### Send Platform-Specific Notification

Send notifications with platform-specific configurations:

```go
// Android-specific notification
androidConfig := &firebase.AndroidConfig{
    Priority: "high",
    Notification: &firebase.AndroidNotification{
        Title:    "Android Alert",
        Body:     "Platform-specific message",
        Icon:     "ic_notification",
        Color:    "#FF0000",
        Sound:    "default",
        ChannelID: "important_channel",
    },
}

// iOS-specific notification
apnsConfig := &firebase.APNSConfig{
    Payload: &firebase.APNSPayload{
        Aps: &firebase.Aps{
            Alert: &firebase.Alert{
                Title: "iOS Alert",
                Body:  "iOS-specific message",
            },
            Badge: intPtr(1),
            Sound: "default",
        },
    },
}

req := firebase.SendToDeviceRequest{
    Message: firebase.Message{
        Token: "device-fcm-token-here",
        Notification: &firebase.Notification{
            Title: "Cross-platform",
            Body:  "Works on all platforms",
        },
        Android: androidConfig,
        APNS:    apnsConfig,
    },
}

response, firebaseErr := client.SendToDevice(ctx, req)
if firebaseErr != nil {
    log.Printf("Error sending platform-specific message: %v", firebaseErr)
    return
}

fmt.Printf("Platform-specific message sent, ID: %s\n", response.Name)
```

## API Reference

### Authentication Methods

#### CreateUser

Creates a new user account.

```go
func (c *Client) CreateUser(ctx context.Context, req CreateUserRequest) (*UserRecord, *FirebaseError)
```

**Request Parameters:**

- `Email` (string, required): User's email address
- `Password` (string, required): User's password
- `PhoneNumber` (string, optional): User's phone number
- `DisplayName` (string, optional): User's display name
- `PhotoURL` (string, optional): User's photo URL
- `Disabled` (bool, optional): Whether the user is disabled

**Response:**

- `UserRecord`: Complete user information including UID, timestamps, etc.

**Timeout:** None (uses Firebase SDK defaults)

#### GetUser

Retrieves user information by UID.

```go
func (c *Client) GetUser(ctx context.Context, req GetUserRequest) (*UserRecord, *FirebaseError)
```

**Request Parameters:**

- `UID` (string, required): User's unique identifier

**Response:**

- `UserRecord`: Complete user information

**Timeout:** None

#### GetUserByEmail

Retrieves user information by email address.

```go
func (c *Client) GetUserByEmail(ctx context.Context, req GetUserByEmailRequest) (*UserRecord, *FirebaseError)
```

**Request Parameters:**

- `Email` (string, required): User's email address

**Response:**

- `UserRecord`: Complete user information

**Timeout:** None

#### GetUserByPhone

Retrieves user information by phone number.

```go
func (c *Client) GetUserByPhone(ctx context.Context, req GetUserByPhoneRequest) (*UserRecord, *FirebaseError)
```

**Request Parameters:**

- `PhoneNumber` (string, required): User's phone number

**Response:**

- `UserRecord`: Complete user information

**Timeout:** None

#### UpdateUser

Updates an existing user's profile information.

```go
func (c *Client) UpdateUser(ctx context.Context, req UpdateUserRequest) (*UserRecord, *FirebaseError)
```

**Request Parameters:**

- `UID` (string, required): User's unique identifier
- `Email` (\*string, optional): New email address
- `PhoneNumber` (\*string, optional): New phone number
- `Password` (\*string, optional): New password
- `DisplayName` (\*string, optional): New display name
- `PhotoURL` (\*string, optional): New photo URL
- `Disabled` (\*bool, optional): Account disabled status
- `CustomClaims` (map[string]interface{}, optional): Custom claims for the user

**Response:**

- `UserRecord`: Updated user information

**Timeout:** None

#### DeleteUser

Deletes a user account permanently.

```go
func (c *Client) DeleteUser(ctx context.Context, req DeleteUserRequest) *FirebaseError
```

**Request Parameters:**

- `UID` (string, required): User's unique identifier

**Response:**

- `nil` on success, `FirebaseError` on failure

**Timeout:** None

#### VerifyIDToken

Verifies and decodes a Firebase ID token.

```go
func (c *Client) VerifyIDToken(ctx context.Context, req VerifyIDTokenRequest) (*Token, *FirebaseError)
```

**Request Parameters:**

- `IDToken` (string, required): Firebase ID token from client

**Response:**

- `Token`: Decoded token information including claims, issuer, subject, etc.

**Timeout:** None

#### CreateCustomToken

Creates a custom token for client-side authentication.

```go
func (c *Client) CreateCustomToken(ctx context.Context, req CreateCustomTokenRequest) (string, *FirebaseError)
```

**Request Parameters:**

- `UID` (string, required): User's unique identifier
- `CustomClaims` (map[string]interface{}, optional): Additional claims to include

**Response:**

- `string`: Custom JWT token

**Timeout:** None

#### RevokeRefreshTokens

Revokes all refresh tokens for a user, forcing re-authentication.

```go
func (c *Client) RevokeRefreshTokens(ctx context.Context, req RevokeRefreshTokensRequest) *FirebaseError
```

**Request Parameters:**

- `UID` (string, required): User's unique identifier

**Response:**

- `nil` on success, `FirebaseError` on failure

**Timeout:** None

### FCM Methods

#### SendToDevice

Sends a push notification or data message to a specific device.

```go
func (c *Client) SendToDevice(ctx context.Context, req SendToDeviceRequest) (*MessageResponse, *FirebaseError)
```

**Request Parameters:**

- `Message.Token` (string, required): FCM registration token of the target device
- `Message.Notification` (\*Notification, optional): Notification payload
- `Message.Data` (map[string]string, optional): Data payload
- `Message.Android` (\*AndroidConfig, optional): Android-specific configuration
- `Message.APNS` (\*APNSConfig, optional): iOS-specific configuration
- `Message.Webpush` (\*WebpushConfig, optional): Web-specific configuration

**Response:**

- `MessageResponse`: Contains the message ID/name

**Timeout:** None

## Error Handling

The package returns `*FirebaseError` for all Firebase-specific errors:

```go
type FirebaseError struct {
    Code    FirebaseErrorCode        `json:"code"`
    Message string                   `json:"message"`
    Details map[string]interface{}   `json:"details,omitempty"`
}

func (e *FirebaseError) Error() string {
    return fmt.Sprintf("Firebase error (%s): %s", e.Code, e.Message)
}
```

### Error Codes

**Authentication Errors:**

- `USER_NOT_FOUND`: User does not exist
- `EMAIL_ALREADY_EXISTS`: Email address already in use
- `PHONE_NUMBER_ALREADY_EXISTS`: Phone number already in use
- `INVALID_EMAIL`: Invalid email address format
- `INVALID_PHONE_NUMBER`: Invalid phone number format
- `WEAK_PASSWORD`: Password does not meet requirements
- `INVALID_ID_TOKEN`: Invalid or expired ID token
- `INSUFFICIENT_PERMISSION`: Insufficient permissions for operation
- `INVALID_ARGUMENT`: Invalid request parameters

**Messaging Errors:**

- `MESSAGING_ERROR`: Generic messaging error
- `INVALID_REGISTRATION`: Invalid FCM registration token
- `UNREGISTERED`: Device token is unregistered
- `INVALID_PACKAGE_NAME`: Invalid package name for FCM

**Generic Errors:**

- `INTERNAL_ERROR`: Internal server error
- `UNKNOWN_ERROR`: Unknown error occurred

**Example Error Handling:**

```go
user, firebaseErr := client.CreateUser(ctx, req)
if firebaseErr != nil {
    switch firebaseErr.Code {
    case firebase.EmailAlreadyExists:
        log.Printf("Email already exists: %s", req.Email)
        // Handle duplicate email
    case firebase.WeakPassword:
        log.Printf("Password too weak for user: %s", req.Email)
        // Handle weak password
    case firebase.InvalidEmail:
        log.Printf("Invalid email format: %s", req.Email)
        // Handle invalid email
    default:
        log.Printf("Unexpected error: %v", firebaseErr)
        // Handle other errors
    }
    return
}
```

## Data Models

### UserRecord

Represents a Firebase user account:

```go
type UserRecord struct {
    UID           string                 `json:"uid"`
    Email         string                 `json:"email"`
    PhoneNumber   string                 `json:"phoneNumber"`
    DisplayName   string                 `json:"displayName"`
    PhotoURL      string                 `json:"photoURL"`
    Disabled      bool                   `json:"disabled"`
    EmailVerified bool                   `json:"emailVerified"`
    CustomClaims  map[string]interface{} `json:"customClaims"`
    CreatedAt     time.Time              `json:"createdAt"`
    LastLoginAt   *time.Time             `json:"lastLoginAt"`
}
```

### Token

Represents a decoded Firebase ID token:

```go
type Token struct {
    Issuer    string                 `json:"iss"`
    Audience  string                 `json:"aud"`
    Subject   string                 `json:"sub"`
    IssuedAt  int64                  `json:"iat"`
    ExpiresAt int64                  `json:"exp"`
    Claims    map[string]interface{} `json:"claims"`
}
```

### Message

Represents an FCM message:

```go
type Message struct {
    Token        string                 `json:"token"`
    Data         map[string]string      `json:"data"`
    Notification *Notification          `json:"notification"`
    Android      *AndroidConfig         `json:"android"`
    APNS         *APNSConfig            `json:"apns"`
    Webpush      *WebpushConfig         `json:"webpush"`
}
```

## Examples

### Complete User Management Flow

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/TixiaOTA/tixia-pkg/providers/firebase"
    "github.com/TixiaOTA/tixia-pkg/logger"
)

func main() {
    // Read service account JSON
    serviceAccountJSON, err := os.ReadFile("path/to/service-account.json")
    if err != nil {
        log.Fatal(err)
    }

    // Create client
    config := firebase.Config{
        ProjectID:          "my-firebase-project",
        ServiceAccountJSON: serviceAccountJSON,
        Logger:             logger.New(),
    }

    client, err := firebase.New(config)
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Create a new user
    createReq := firebase.CreateUserRequest{
        Email:       "john.doe@example.com",
        Password:    "securePassword123!",
        DisplayName: "John Doe",
        PhoneNumber: "+1234567890",
    }

    user, firebaseErr := client.CreateUser(ctx, createReq)
    if firebaseErr != nil {
        log.Fatal(firebaseErr)
    }

    fmt.Printf("Created user: %s (UID: %s)\n", user.DisplayName, user.UID)

    // Get user by email
    getReq := firebase.GetUserByEmailRequest{
        Email: "john.doe@example.com",
    }

    retrievedUser, firebaseErr := client.GetUserByEmail(ctx, getReq)
    if firebaseErr != nil {
        log.Fatal(firebaseErr)
    }

    fmt.Printf("Retrieved user: %s\n", retrievedUser.DisplayName)

    // Update user
    updateReq := firebase.UpdateUserRequest{
        UID:         user.UID,
        DisplayName: stringPtr("John Smith"),
        CustomClaims: map[string]interface{}{
            "role": "user",
            "verified": true,
        },
    }

    updatedUser, firebaseErr := client.UpdateUser(ctx, updateReq)
    if firebaseErr != nil {
        log.Fatal(firebaseErr)
    }

    fmt.Printf("Updated user: %s\n", updatedUser.DisplayName)

    // Send welcome notification
    msgReq := firebase.SendToDeviceRequest{
        Message: firebase.Message{
            Token: "user-device-token-here", // Would come from user device
            Notification: &firebase.Notification{
                Title: "Welcome!",
                Body:  "Thanks for joining our app, " + updatedUser.DisplayName,
            },
            Data: map[string]string{
                "user_id": updatedUser.UID,
                "type":    "welcome",
            },
        },
    }

    msgResp, firebaseErr := client.SendToDevice(ctx, msgReq)
    if firebaseErr != nil {
        log.Printf("Failed to send notification: %v", firebaseErr)
    } else {
        fmt.Printf("Welcome notification sent: %s\n", msgResp.Name)
    }
}

// Helper functions
func stringPtr(s string) *string {
    return &s
}

func boolPtr(b bool) *bool {
    return &b
}
```

## License

This package is part of the TixiaOTA project.
