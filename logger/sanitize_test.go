package logger

import (
	"encoding/json"
	"testing"
)

// Example struct with sensitive fields
type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password" log:"redact"`
	Email    string `json:"email"`
	Token    string `json:"token" log:"redact"`
}

// Example nested struct
type LoginRequest struct {
	Credentials UserCredentials `json:"credentials"`
	IPAddress   string          `json:"ip_address"`
	UserAgent   string          `json:"user_agent"`
}

// Example struct with slices
type UserList struct {
	Users []UserCredentials `json:"users"`
	Count int               `json:"count"`
}

// Example struct with pointer fields
type APIConfig struct {
	APIKey    *string `json:"api_key" log:"redact"`
	SecretKey *string `json:"secret_key" log:"redact"`
	Endpoint  string  `json:"endpoint"`
}

func TestSanitize_BasicStruct(t *testing.T) {
	creds := UserCredentials{
		Username: "john_doe",
		Password: "secret123",
		Email:    "john@example.com",
		Token:    "abc123xyz",
	}

	sanitized, err := Sanitize(creds)
	if err != nil {
		t.Fatalf("Sanitize failed: %v", err)
	}

	jsonStr, _ := json.MarshalIndent(sanitized, "", "  ")
	t.Logf("Input: %+v", creds)
	t.Logf("Sanitized result:\n%s", string(jsonStr))

	// Check that password is redacted
	if sanitized["password"] != "[REDACTED]" {
		t.Errorf("Expected password to be redacted, got: %v", sanitized["password"])
	}

	// Check that token is redacted
	if sanitized["token"] != "[REDACTED]" {
		t.Errorf("Expected token to be redacted, got: %v", sanitized["token"])
	}

	// Check that non-sensitive fields are preserved
	if sanitized["username"] != "john_doe" {
		t.Errorf("Expected username to be preserved, got: %v", sanitized["username"])
	}

	if sanitized["email"] != "john@example.com" {
		t.Errorf("Expected email to be preserved, got: %v", sanitized["email"])
	}
}

func TestSanitize_NestedStruct(t *testing.T) {
	request := LoginRequest{
		Credentials: UserCredentials{
			Username: "jane_doe",
			Password: "password123",
			Email:    "jane@example.com",
			Token:    "token456",
		},
		IPAddress: "192.168.1.1",
		UserAgent: "Mozilla/5.0",
	}

	sanitized, err := Sanitize(request)
	if err != nil {
		t.Fatalf("Sanitize failed: %v", err)
	}

	jsonStr, _ := json.MarshalIndent(sanitized, "", "  ")
	t.Logf("Input: %+v", request)
	t.Logf("Sanitized result:\n%s", string(jsonStr))

	// Check nested credentials
	creds, ok := sanitized["credentials"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected credentials to be a map, got: %T", sanitized["credentials"])
	}

	if creds["password"] != "[REDACTED]" {
		t.Errorf("Expected nested password to be redacted, got: %v", creds["password"])
	}

	if creds["token"] != "[REDACTED]" {
		t.Errorf("Expected nested token to be redacted, got: %v", creds["token"])
	}

	// Check top-level fields
	if sanitized["ip_address"] != "192.168.1.1" {
		t.Errorf("Expected ip_address to be preserved, got: %v", sanitized["ip_address"])
	}
}

func TestSanitize_WithPointer(t *testing.T) {
	apiKey := "secret-api-key"
	secretKey := "super-secret"
	endpoint := "https://api.example.com"

	config := APIConfig{
		APIKey:    &apiKey,
		SecretKey: &secretKey,
		Endpoint:  endpoint,
	}

	sanitized, err := Sanitize(config)
	if err != nil {
		t.Fatalf("Sanitize failed: %v", err)
	}

	jsonStr, _ := json.MarshalIndent(sanitized, "", "  ")
	t.Logf("Input: APIKey=%s, SecretKey=%s, Endpoint=%s", *config.APIKey, *config.SecretKey, config.Endpoint)
	t.Logf("Sanitized result:\n%s", string(jsonStr))

	// Check that pointer fields are redacted
	if sanitized["api_key"] != "[REDACTED]" {
		t.Errorf("Expected api_key to be redacted, got: %v", sanitized["api_key"])
	}

	if sanitized["secret_key"] != "[REDACTED]" {
		t.Errorf("Expected secret_key to be redacted, got: %v", sanitized["secret_key"])
	}

	// Check that non-sensitive field is preserved
	if sanitized["endpoint"] != endpoint {
		t.Errorf("Expected endpoint to be preserved, got: %v", sanitized["endpoint"])
	}
}

func TestSanitize_WithSlice(t *testing.T) {
	userList := UserList{
		Users: []UserCredentials{
			{
				Username: "user1",
				Password: "pass1",
				Email:    "user1@example.com",
				Token:    "token1",
			},
			{
				Username: "user2",
				Password: "pass2",
				Email:    "user2@example.com",
				Token:    "token2",
			},
		},
		Count: 2,
	}

	sanitized, err := Sanitize(userList)
	if err != nil {
		t.Fatalf("Sanitize failed: %v", err)
	}

	jsonStr, _ := json.MarshalIndent(sanitized, "", "  ")
	t.Logf("Input: %d users", len(userList.Users))
	t.Logf("Sanitized result:\n%s", string(jsonStr))

	users, ok := sanitized["users"].([]interface{})
	if !ok {
		t.Fatalf("Expected users to be a slice, got: %T", sanitized["users"])
	}

	if len(users) != 2 {
		t.Errorf("Expected 2 users, got: %d", len(users))
	}

	// Check first user
	user1, ok := users[0].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected user1 to be a map, got: %T", users[0])
	}

	if user1["password"] != "[REDACTED]" {
		t.Errorf("Expected user1 password to be redacted, got: %v", user1["password"])
	}
}

func TestSanitize_WithPointerInput(t *testing.T) {
	creds := &UserCredentials{
		Username: "pointer_user",
		Password: "pointer_pass",
		Email:    "pointer@example.com",
		Token:    "pointer_token",
	}

	sanitized, err := Sanitize(creds)
	if err != nil {
		t.Fatalf("Sanitize failed: %v", err)
	}

	jsonStr, _ := json.MarshalIndent(sanitized, "", "  ")
	t.Logf("Input (pointer): %+v", *creds)
	t.Logf("Sanitized result:\n%s", string(jsonStr))

	if sanitized["password"] != "[REDACTED]" {
		t.Errorf("Expected password to be redacted, got: %v", sanitized["password"])
	}
}

func TestSanitize_NilPointer(t *testing.T) {
	var creds *UserCredentials = nil

	sanitized, err := Sanitize(creds)
	if err != nil {
		t.Fatalf("Sanitize failed: %v", err)
	}

	if sanitized != nil {
		t.Errorf("Expected nil result for nil pointer input, got: %v", sanitized)
	}
}

func TestSanitizeToJSON(t *testing.T) {
	creds := UserCredentials{
		Username: "json_user",
		Password: "json_pass",
		Email:    "json@example.com",
		Token:    "json_token",
	}

	jsonStr, err := SanitizeToJSON(creds)
	if err != nil {
		t.Fatalf("SanitizeToJSON failed: %v", err)
	}

	t.Logf("Input: %+v", creds)
	t.Logf("Sanitized JSON result:\n%s", jsonStr)

	// Verify it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("SanitizeToJSON returned invalid JSON: %v", err)
	}

	// Verify redaction
	if result["password"] != "[REDACTED]" {
		t.Errorf("Expected password to be redacted in JSON, got: %v", result["password"])
	}
}

func TestSanitize_InvalidInput(t *testing.T) {
	invalidInput := "not a struct"

	t.Logf("Testing invalid input: %v", invalidInput)
	_, err := Sanitize(invalidInput)
	if err == nil {
		t.Error("Expected error for non-struct input, got nil")
	} else {
		t.Logf("Expected error received: %v", err)
	}
}
