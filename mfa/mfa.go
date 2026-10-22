package mfa

import (
	"time"

	"github.com/TixiaOTA/tixia-pkg/logger"
)

// Manager provides TOTP enrollment, verification, and recovery code operations.
type Manager struct {
	config Config
	logger *logger.Logger
}

// New creates a new MFA manager.
func New(config Config, logger *logger.Logger) (*Manager, error) {
	config.applyDefaults()

	if err := config.Validate(); err != nil {
		return nil, err
	}

	if logger == nil {
		return nil, ErrNilLogger
	}

	return &Manager{
		config: config,
		logger: logger,
	}, nil
}

func (m *Manager) now() time.Time {
	return m.config.Now()
}

func (m *Manager) replayTTL() time.Duration {
	periods := (2 * m.config.skew()) + 2
	return time.Duration(periods*m.config.Period) * time.Second
}
