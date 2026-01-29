package firebase

import (
	"context"
	"encoding/json"
	"errors"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/messaging"
	"github.com/TixiaOTA/tixia-pkg/logger"
	"google.golang.org/api/option"
)

// Config holds the configuration for the Firebase client
type Config struct {
	ProjectID          string         // Required: Firebase project ID
	ServiceAccountJSON []byte         // Required: Service account key JSON as bytes
	Logger             *logger.Logger // Required: tixia-pkg logger
}

// Client represents the Firebase client that provides authentication and messaging services
type Client struct {
	config     Config
	app        *firebase.App
	authClient *auth.Client
	msgClient  *messaging.Client
	logger     *logger.Logger
}

// New creates a new Firebase client instance with service account authentication
func New(config Config) (*Client, error) {
	// Validate config
	if config.ProjectID == "" {
		return nil, errors.New("project ID is required")
	}
	if len(config.ServiceAccountJSON) == 0 {
		return nil, errors.New("service account JSON is required")
	}
	if config.Logger == nil {
		return nil, errors.New("logger is required")
	}

	// Validate service account JSON
	var sa map[string]interface{}
	if err := json.Unmarshal(config.ServiceAccountJSON, &sa); err != nil {
		return nil, &FirebaseError{
			Code:    InvalidArgument,
			Message: "invalid service account JSON",
			Details: map[string]interface{}{
				"originalError": err.Error(),
			},
		}
	}

	// Verify required fields in service account
	requiredFields := []string{"type", "project_id", "private_key_id", "private_key", "client_email", "client_id"}
	for _, field := range requiredFields {
		if _, exists := sa[field]; !exists {
			return nil, &FirebaseError{
				Code:    InvalidArgument,
				Message: "service account JSON missing required field: " + field,
			}
		}
	}

	// Initialize Firebase app with service account
	firebaseConfig := &firebase.Config{
		ProjectID: config.ProjectID,
	}

	app, err := firebase.NewApp(context.Background(), firebaseConfig, option.WithCredentialsJSON(config.ServiceAccountJSON))
	if err != nil {
		return nil, &FirebaseError{
			Code:    InternalError,
			Message: "failed to initialize Firebase app",
			Details: map[string]interface{}{
				"originalError": err.Error(),
			},
		}
	}

	// Initialize Firebase Auth client
	authClient, err := app.Auth(context.Background())
	if err != nil {
		return nil, &FirebaseError{
			Code:    InternalError,
			Message: "failed to initialize auth client",
			Details: map[string]interface{}{
				"originalError": err.Error(),
			},
		}
	}

	// Initialize Firebase Messaging client
	msgClient, err := app.Messaging(context.Background())
	if err != nil {
		return nil, &FirebaseError{
			Code:    InternalError,
			Message: "failed to initialize messaging client",
			Details: map[string]interface{}{
				"originalError": err.Error(),
			},
		}
	}

	return &Client{
		config:     config,
		app:        app,
		authClient: authClient,
		msgClient:  msgClient,
		logger:     config.Logger,
	}, nil
}
