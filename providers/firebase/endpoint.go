package firebase

import (
	"context"
	"maps"

	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/messaging"
	"github.com/TixiaOTA/tixia-pkg/logger"
	otelpkg "github.com/TixiaOTA/tixia-pkg/otel"
	"go.uber.org/zap"
)

// --------------------------------- AUTHENTICATION METHODS ---------------------------------

// CreateUser creates a new user account
func (c *Client) CreateUser(ctx context.Context, req CreateUserRequest) (*UserRecord, *FirebaseError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "firebase", "CreateUser")
	defer span.End()

	// Sanitize request for logging
	sanitizedReq, _ := logger.Sanitize(req)
	c.logger.WithTrace(ctx, currentSpanId).Info("Starting CreateUser request", zap.Any("request", sanitizedReq))

	// Create user params
	params := (&auth.UserToCreate{}).
		Email(req.Email).
		Password(req.Password)

	if req.PhoneNumber != "" {
		params = params.PhoneNumber(req.PhoneNumber)
	}
	if req.DisplayName != "" {
		params = params.DisplayName(req.DisplayName)
	}
	if req.PhotoURL != "" {
		params = params.PhotoURL(req.PhotoURL)
	}
	if req.Disabled {
		params = params.Disabled(req.Disabled)
	}

	// Create user
	user, err := c.authClient.CreateUser(ctx, params)
	if err != nil {
		firebaseErr := mapAuthError(err)
		c.logger.WithTrace(ctx, currentSpanId).Error("CreateUser failed", zap.Error(firebaseErr))
		return nil, firebaseErr
	}

	userRecord := convertUserRecord(user)
	c.logger.WithTrace(ctx, currentSpanId).Info("CreateUser completed successfully", zap.String("uid", userRecord.UID))
	return userRecord, nil
}

// GetUser retrieves a user by UID
func (c *Client) GetUser(ctx context.Context, req GetUserRequest) (*UserRecord, *FirebaseError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "firebase", "GetUser")
	defer span.End()

	c.logger.WithTrace(ctx, currentSpanId).Info("Starting GetUser request", zap.String("uid", req.UID))

	user, err := c.authClient.GetUser(ctx, req.UID)
	if err != nil {
		firebaseErr := mapAuthError(err)
		c.logger.WithTrace(ctx, currentSpanId).Error("GetUser failed", zap.Error(firebaseErr))
		return nil, firebaseErr
	}

	userRecord := convertUserRecord(user)
	c.logger.WithTrace(ctx, currentSpanId).Info("GetUser completed successfully", zap.String("uid", userRecord.UID))
	return userRecord, nil
}

// GetUserByEmail retrieves a user by email
func (c *Client) GetUserByEmail(ctx context.Context, req GetUserByEmailRequest) (*UserRecord, *FirebaseError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "firebase", "GetUserByEmail")
	defer span.End()

	c.logger.WithTrace(ctx, currentSpanId).Info("Starting GetUserByEmail request", zap.String("email", req.Email))

	user, err := c.authClient.GetUserByEmail(ctx, req.Email)
	if err != nil {
		firebaseErr := mapAuthError(err)
		c.logger.WithTrace(ctx, currentSpanId).Error("GetUserByEmail failed", zap.Error(firebaseErr))
		return nil, firebaseErr
	}

	userRecord := convertUserRecord(user)
	c.logger.WithTrace(ctx, currentSpanId).Info("GetUserByEmail completed successfully", zap.String("uid", userRecord.UID))
	return userRecord, nil
}

// GetUserByPhone retrieves a user by phone number
func (c *Client) GetUserByPhone(ctx context.Context, req GetUserByPhoneRequest) (*UserRecord, *FirebaseError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "firebase", "GetUserByPhone")
	defer span.End()

	c.logger.WithTrace(ctx, currentSpanId).Info("Starting GetUserByPhone request", zap.String("phoneNumber", req.PhoneNumber))

	user, err := c.authClient.GetUserByPhoneNumber(ctx, req.PhoneNumber)
	if err != nil {
		firebaseErr := mapAuthError(err)
		c.logger.WithTrace(ctx, currentSpanId).Error("GetUserByPhone failed", zap.Error(firebaseErr))
		return nil, firebaseErr
	}

	userRecord := convertUserRecord(user)
	c.logger.WithTrace(ctx, currentSpanId).Info("GetUserByPhone completed successfully", zap.String("uid", userRecord.UID))
	return userRecord, nil
}

// UpdateUser updates an existing user
func (c *Client) UpdateUser(ctx context.Context, req UpdateUserRequest) (*UserRecord, *FirebaseError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "firebase", "UpdateUser")
	defer span.End()

	// Sanitize request for logging
	sanitizedReq, _ := logger.Sanitize(req)
	c.logger.WithTrace(ctx, currentSpanId).Info("Starting UpdateUser request", zap.Any("request", sanitizedReq))

	// Create update params
	params := (&auth.UserToUpdate{})

	if req.Email != nil {
		params = params.Email(*req.Email)
	}
	if req.PhoneNumber != nil {
		params = params.PhoneNumber(*req.PhoneNumber)
	}
	if req.Password != nil {
		params = params.Password(*req.Password)
	}
	if req.DisplayName != nil {
		params = params.DisplayName(*req.DisplayName)
	}
	if req.PhotoURL != nil {
		params = params.PhotoURL(*req.PhotoURL)
	}
	if req.Disabled != nil {
		params = params.Disabled(*req.Disabled)
	}
	if req.CustomClaims != nil {
		params = params.CustomClaims(req.CustomClaims)
	}

	// Update user
	user, err := c.authClient.UpdateUser(ctx, req.UID, params)
	if err != nil {
		firebaseErr := mapAuthError(err)
		c.logger.WithTrace(ctx, currentSpanId).Error("UpdateUser failed", zap.Error(firebaseErr))
		return nil, firebaseErr
	}

	userRecord := convertUserRecord(user)
	c.logger.WithTrace(ctx, currentSpanId).Info("UpdateUser completed successfully", zap.String("uid", userRecord.UID))
	return userRecord, nil
}

// DeleteUser deletes a user account
func (c *Client) DeleteUser(ctx context.Context, req DeleteUserRequest) *FirebaseError {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "firebase", "DeleteUser")
	defer span.End()

	c.logger.WithTrace(ctx, currentSpanId).Info("Starting DeleteUser request", zap.String("uid", req.UID))

	err := c.authClient.DeleteUser(ctx, req.UID)
	if err != nil {
		firebaseErr := mapAuthError(err)
		c.logger.WithTrace(ctx, currentSpanId).Error("DeleteUser failed", zap.Error(firebaseErr))
		return firebaseErr
	}

	c.logger.WithTrace(ctx, currentSpanId).Info("DeleteUser completed successfully", zap.String("uid", req.UID))
	return nil
}

// VerifyIDToken verifies an ID token and returns the decoded token
func (c *Client) VerifyIDToken(ctx context.Context, req VerifyIDTokenRequest) (*Token, *FirebaseError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "firebase", "VerifyIDToken")
	defer span.End()

	c.logger.WithTrace(ctx, currentSpanId).Info("Starting VerifyIDToken request")

	token, err := c.authClient.VerifyIDToken(ctx, req.IDToken)
	if err != nil {
		firebaseErr := mapAuthError(err)
		c.logger.WithTrace(ctx, currentSpanId).Error("VerifyIDToken failed", zap.Error(firebaseErr))
		return nil, firebaseErr
	}

	tokenRecord := convertToken(token)
	c.logger.WithTrace(ctx, currentSpanId).Info("VerifyIDToken completed successfully", zap.String("subject", tokenRecord.Subject))
	return tokenRecord, nil
}

// CreateCustomToken creates a custom token for the given user
func (c *Client) CreateCustomToken(ctx context.Context, req CreateCustomTokenRequest) (string, *FirebaseError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "firebase", "CreateCustomToken")
	defer span.End()

	c.logger.WithTrace(ctx, currentSpanId).Info("Starting CreateCustomToken request", zap.String("uid", req.UID))

	token, err := c.authClient.CustomToken(ctx, req.UID)
	if err != nil {
		firebaseErr := mapAuthError(err)
		c.logger.WithTrace(ctx, currentSpanId).Error("CreateCustomToken failed", zap.Error(firebaseErr))
		return "", firebaseErr
	}

	c.logger.WithTrace(ctx, currentSpanId).Info("CreateCustomToken completed successfully", zap.String("uid", req.UID))
	return token, nil
}

// RevokeRefreshTokens revokes all refresh tokens for a user
func (c *Client) RevokeRefreshTokens(ctx context.Context, req RevokeRefreshTokensRequest) *FirebaseError {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "firebase", "RevokeRefreshTokens")
	defer span.End()

	c.logger.WithTrace(ctx, currentSpanId).Info("Starting RevokeRefreshTokens request", zap.String("uid", req.UID))

	err := c.authClient.RevokeRefreshTokens(ctx, req.UID)
	if err != nil {
		firebaseErr := mapAuthError(err)
		c.logger.WithTrace(ctx, currentSpanId).Error("RevokeRefreshTokens failed", zap.Error(firebaseErr))
		return firebaseErr
	}

	c.logger.WithTrace(ctx, currentSpanId).Info("RevokeRefreshTokens completed successfully", zap.String("uid", req.UID))
	return nil
}

// --------------------------------- FCM METHODS ---------------------------------

// SendToDevice sends a message to a specific device
func (c *Client) SendToDevice(ctx context.Context, req SendToDeviceRequest) (*MessageResponse, *FirebaseError) {
	ctx, span, currentSpanId := otelpkg.StartSpanWithID(ctx, "provider", "firebase", "SendToDevice")
	defer span.End()

	if len(req.Message.Tokens) == 0 {
		return nil, &FirebaseError{
			Code:    InvalidArgument,
			Message: "no tokens provided",
		}
	}

	// Sanitize request for logging (remove sensitive data)
	sanitizedReq, _ := logger.Sanitize(req)
	c.logger.WithTrace(ctx, currentSpanId).Info("Starting SendToDevice request", zap.Any("request", sanitizedReq))

	// Convert our message to Firebase SDK message
	fcmMessage := &messaging.Message{}

	// Add notification if present
	if req.Message.Notification != nil {
		fcmMessage.Notification = &messaging.Notification{
			Title:    req.Message.Notification.Title,
			Body:     req.Message.Notification.Body,
			ImageURL: req.Message.Notification.Image,
		}
	}

	// Add Android config if present
	if req.Message.Android != nil {
		fcmMessage.Android = &messaging.AndroidConfig{
			CollapseKey:           req.Message.Android.CollapseKey,
			Priority:              req.Message.Android.Priority,
			TTL:                   req.Message.Android.TTL,
			RestrictedPackageName: req.Message.Android.RestrictedPackageName,
			Data:                  req.Message.Android.Data,
		}

		if req.Message.Android.Notification != nil {
			fcmMessage.Android.Notification = &messaging.AndroidNotification{
				Title:        req.Message.Android.Notification.Title,
				Body:         req.Message.Android.Notification.Body,
				Icon:         req.Message.Android.Notification.Icon,
				Color:        req.Message.Android.Notification.Color,
				Sound:        req.Message.Android.Notification.Sound,
				Tag:          req.Message.Android.Notification.Tag,
				ClickAction:  req.Message.Android.Notification.ClickAction,
				BodyLocKey:   req.Message.Android.Notification.BodyLocKey,
				BodyLocArgs:  req.Message.Android.Notification.BodyLocArgs,
				TitleLocKey:  req.Message.Android.Notification.TitleLocKey,
				TitleLocArgs: req.Message.Android.Notification.TitleLocArgs,
				ChannelID:    req.Message.Android.Notification.ChannelID,
				ImageURL:     req.Message.Android.Notification.Image,
			}
		}
	}

	// Add APNS config if present
	if req.Message.APNS != nil {
		fcmMessage.APNS = &messaging.APNSConfig{
			Headers: req.Message.APNS.Headers,
		}

		if req.Message.APNS.Payload != nil {
			fcmMessage.APNS.Payload = &messaging.APNSPayload{
				CustomData: req.Message.APNS.Payload.CustomData,
			}

			if req.Message.APNS.Payload.Aps != nil {
				aps := &messaging.Aps{
					Badge:    req.Message.APNS.Payload.Aps.Badge,
					Sound:    req.Message.APNS.Payload.Aps.Sound,
					Category: req.Message.APNS.Payload.Aps.Category,
					ThreadID: req.Message.APNS.Payload.Aps.ThreadID,
				}

				if req.Message.APNS.Payload.Aps.ContentAvailable != nil {
					aps.ContentAvailable = *req.Message.APNS.Payload.Aps.ContentAvailable != 0
				}

				fcmMessage.APNS.Payload.Aps = aps
			}
		}
	}

	// Add Webpush config if present
	if req.Message.Webpush != nil {
		fcmMessage.Webpush = &messaging.WebpushConfig{
			Headers: req.Message.Webpush.Headers,
			Data:    req.Message.Webpush.Data,
		}

		if req.Message.Webpush.Notification != nil {
			fcmMessage.Webpush.Notification = &messaging.WebpushNotification{
				Title: req.Message.Webpush.Notification.Title,
				Body:  req.Message.Webpush.Notification.Body,
				Icon:  req.Message.Webpush.Notification.Icon,
				Image: req.Message.Webpush.Notification.Image,
				Badge: req.Message.Webpush.Notification.Badge,
			}
		}
	}

	if req.Message.Data != nil {
		if fcmMessage.Data == nil {
			fcmMessage.Data = make(map[string]string)
		}
		maps.Copy(fcmMessage.Data, req.Message.Data)
	}

	// Send messages to all device tokens in a single batch request
	messages := make([]*messaging.Message, 0, len(req.Message.Tokens))
	for _, token := range req.Message.Tokens {
		msg := &messaging.Message{
			Token:        token,
			Data:         fcmMessage.Data,
			Notification: fcmMessage.Notification,
			Android:      fcmMessage.Android,
			APNS:         fcmMessage.APNS,
			Webpush:      fcmMessage.Webpush,
		}
		messages = append(messages, msg)
	}

	batchResp, err := c.msgClient.SendEach(ctx, messages)
	if err != nil {
		firebaseErr := mapMessagingError(err)
		c.logger.WithTrace(ctx, currentSpanId).Error("SendToDevice batch send failed", zap.Error(firebaseErr))
		return nil, firebaseErr
	}

	// Collect successful message IDs, log failures
	messageIds := make([]string, 0, batchResp.SuccessCount)
	for i, resp := range batchResp.Responses {
		if resp.Success {
			messageIds = append(messageIds, resp.MessageID)
		} else {
			c.logger.WithTrace(ctx, currentSpanId).Error("SendToDevice partial failure",
				zap.Int("tokenIndex", i),
				zap.String("token", req.Message.Tokens[i]),
				zap.Error(resp.Error))
		}
	}

	messageResponse := &MessageResponse{
		MessageIDs: messageIds,
	}

	c.logger.WithTrace(ctx, currentSpanId).Info("SendToDevice completed successfully", zap.Any("messageIds", messageIds))
	return messageResponse, nil
}
