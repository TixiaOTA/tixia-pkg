package firebase

import "time"

// --------------------------------- AUTHENTICATION TYPES ---------------------------------

// UserRecord represents a Firebase user account
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

// Token represents a decoded Firebase ID token
type Token struct {
	Issuer    string                 `json:"iss"`
	Audience  string                 `json:"aud"`
	Subject   string                 `json:"sub"`
	IssuedAt  int64                  `json:"iat"`
	ExpiresAt int64                  `json:"exp"`
	Claims    map[string]interface{} `json:"claims"`
}

// --------------------------------- AUTHENTICATION REQUESTS ---------------------------------

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName"`
	PhotoURL    string `json:"photoURL"`
	Disabled    bool   `json:"disabled"`
}

// GetUserRequest represents the request to get a user by UID
type GetUserRequest struct {
	UID string `json:"uid"`
}

// GetUserByEmailRequest represents the request to get a user by email
type GetUserByEmailRequest struct {
	Email string `json:"email"`
}

// GetUserByPhoneRequest represents the request to get a user by phone number
type GetUserByPhoneRequest struct {
	PhoneNumber string `json:"phoneNumber"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	UID          string                 `json:"uid"`
	Email        *string                `json:"email,omitempty"`
	PhoneNumber  *string                `json:"phoneNumber,omitempty"`
	Password     *string                `json:"password,omitempty"`
	DisplayName  *string                `json:"displayName,omitempty"`
	PhotoURL     *string                `json:"photoURL,omitempty"`
	Disabled     *bool                  `json:"disabled,omitempty"`
	CustomClaims map[string]interface{} `json:"customClaims,omitempty"`
}

// DeleteUserRequest represents the request to delete a user
type DeleteUserRequest struct {
	UID string `json:"uid"`
}

// VerifyIDTokenRequest represents the request to verify an ID token
type VerifyIDTokenRequest struct {
	IDToken string `json:"idToken"`
}

// CreateCustomTokenRequest represents the request to create a custom token
type CreateCustomTokenRequest struct {
	UID          string                 `json:"uid"`
	CustomClaims map[string]interface{} `json:"customClaims,omitempty"`
}

// RevokeRefreshTokensRequest represents the request to revoke refresh tokens
type RevokeRefreshTokensRequest struct {
	UID string `json:"uid"`
}

// --------------------------------- FCM TYPES ---------------------------------

// Message represents an FCM message
type Message struct {
	Tokens       []string          `json:"tokens"`       // Devices token
	Data         map[string]string `json:"data"`         // Data payload
	Notification *Notification     `json:"notification"` // Notification payload
	Android      *AndroidConfig    `json:"android"`      // Android-specific config
	APNS         *APNSConfig       `json:"apns"`         // iOS-specific config
	Webpush      *WebpushConfig    `json:"webpush"`      // Web-specific config
}

// Notification represents the notification payload
type Notification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Image string `json:"image"`
}

// AndroidConfig represents Android-specific options for messages sent through FCM
type AndroidConfig struct {
	CollapseKey           string               `json:"collapseKey"`
	Priority              string               `json:"priority"` // "normal" or "high"
	TTL                   *time.Duration       `json:"ttl"`
	RestrictedPackageName string               `json:"restrictedPackageName"`
	Data                  map[string]string    `json:"data"`
	Notification          *AndroidNotification `json:"notification"`
}

// AndroidNotification represents Android-specific notification options
type AndroidNotification struct {
	Title        string   `json:"title"`
	Body         string   `json:"body"`
	Icon         string   `json:"icon"`
	Color        string   `json:"color"`
	Sound        string   `json:"sound"`
	Tag          string   `json:"tag"`
	ClickAction  string   `json:"clickAction"`
	BodyLocKey   string   `json:"bodyLocKey"`
	BodyLocArgs  []string `json:"bodyLocArgs"`
	TitleLocKey  string   `json:"titleLocKey"`
	TitleLocArgs []string `json:"titleLocArgs"`
	ChannelID    string   `json:"channelId"`
	Image        string   `json:"image"`
}

// APNSConfig represents Apple Push Notification Service specific options
type APNSConfig struct {
	Headers map[string]string `json:"headers"`
	Payload *APNSPayload      `json:"payload"`
}

// APNSPayload represents the APNS payload
type APNSPayload struct {
	Aps        *Aps                   `json:"aps"`
	CustomData map[string]interface{} `json:"customData,omitempty"`
}

// Aps represents the APS dictionary
type Aps struct {
	Alert            *Alert `json:"alert,omitempty"`
	Badge            *int   `json:"badge,omitempty"`
	Sound            string `json:"sound,omitempty"`
	ContentAvailable *int   `json:"contentAvailable,omitempty"`
	Category         string `json:"category,omitempty"`
	ThreadID         string `json:"threadId,omitempty"`
}

// Alert represents the alert payload
type Alert struct {
	Title        string   `json:"title,omitempty"`
	Body         string   `json:"body,omitempty"`
	TitleLocKey  string   `json:"titleLocKey,omitempty"`
	TitleLocArgs []string `json:"titleLocArgs,omitempty"`
	ActionLocKey string   `json:"actionLocKey,omitempty"`
	LocKey       string   `json:"locKey,omitempty"`
	LocArgs      []string `json:"locArgs,omitempty"`
	LaunchImage  string   `json:"launchImage,omitempty"`
}

// WebpushConfig represents Webpush specific options
type WebpushConfig struct {
	Headers      map[string]string    `json:"headers"`
	Data         map[string]string    `json:"data"`
	Notification *WebpushNotification `json:"notification"`
}

// WebpushNotification represents Webpush notification options
type WebpushNotification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Icon  string `json:"icon"`
	Image string `json:"image"`
	Badge string `json:"badge"`
}

// MessageResponse represents the response from FCM
type MessageResponse struct {
	MessageIDs []string `json:"messageIds"` // Message IDs
}

// --------------------------------- FCM REQUESTS ---------------------------------

// SendToDeviceRequest represents the request to send a message to a device
type SendToDeviceRequest struct {
	Message Message `json:"message"`
}
