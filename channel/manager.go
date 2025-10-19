package channel

import (
	"fmt"
	"sync"
	"time"

	"github.com/TixiaOTA/tixia-backend-master/pkg/logger"
	"go.uber.org/zap"
)

// Manager is a generic channel manager for handling real-time event broadcasting
type Manager[T any] struct {
	clients    map[string]map[string]chan T // topic -> clientID -> channel
	mutex      sync.RWMutex
	bufferSize int
	logger     *logger.Logger
}

// NewManager creates a new channel manager with the specified buffer size
func NewManager[T any](bufferSize int, logger *logger.Logger) *Manager[T] {
	if bufferSize <= 0 {
		bufferSize = 100 // default buffer size
	}

	return &Manager[T]{
		clients:    make(map[string]map[string]chan T),
		bufferSize: bufferSize,
		logger:     logger,
	}
}

// RegisterClient registers a new client for a topic and returns the client ID
func (m *Manager[T]) RegisterClient(topic string, clientID string) (string, chan T) {
	if clientID == "" {
		clientID = generateClientID()
	}

	clientChan := make(chan T, m.bufferSize)

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Check if client ID already exists for this topic (collision detection)
	if existingClients, exists := m.clients[topic]; exists {
		if _, clientExists := existingClients[clientID]; clientExists {
			m.logger.Warn("Client ID collision detected!",
				zap.String("topic", topic),
				zap.String("colliding_client_id", clientID),
				zap.Int("existing_clients_count", len(existingClients)),
			)
			// Generate a new unique client ID to avoid collision
			clientID = generateClientID()
			m.logger.Info("Generated new client ID to resolve collision",
				zap.String("topic", topic),
				zap.String("new_client_id", clientID),
			)
		}
	}

	if m.clients[topic] == nil {
		m.clients[topic] = make(map[string]chan T)
	}
	m.clients[topic][clientID] = clientChan

	m.logger.Info("Registered client",
		zap.String("topic", topic),
		zap.String("client_id", clientID),
		zap.Int("total_clients_for_topic", len(m.clients[topic])),
	)

	return clientID, clientChan
}

// UnregisterClient removes a client from a topic
func (m *Manager[T]) UnregisterClient(topic, clientID string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if clients, exists := m.clients[topic]; exists {
		if clientChan, exists := clients[clientID]; exists {
			close(clientChan)
			delete(clients, clientID)

			m.logger.Info("Unregistered client from topic",
				zap.String("topic", topic),
				zap.String("client_id", clientID),
				zap.Int("remaining_clients", len(clients)),
			)

			// Clean up empty topic entries
			if len(clients) == 0 {
				delete(m.clients, topic)
			}
		} else {
			m.logger.Warn("Attempted to unregister non-existent client",
				zap.String("topic", topic),
				zap.String("client_id", clientID),
			)
		}
	} else {
		m.logger.Warn("Attempted to unregister client from non-existent topic",
			zap.String("topic", topic),
			zap.String("client_id", clientID),
		)
	}
}

// Broadcast sends an event to all clients subscribed to a topic
func (m *Manager[T]) Broadcast(topic string, event T) {
	m.mutex.RLock()
	clients, exists := m.clients[topic]
	clientCount := 0
	if exists {
		clientCount = len(clients)
	}
	m.mutex.RUnlock()

	if !exists {
		m.logger.Debug("No clients found for topic",
			zap.String("topic", topic),
		)
		return
	}

	m.logger.Debug("Broadcasting event to clients",
		zap.String("topic", topic),
		zap.Int("client_count", clientCount),
	)

	// Broadcast to all connected clients for this topic
	successCount := 0
	for clientID, clientChan := range clients {
		select {
		case clientChan <- event:
			// Event sent successfully
			successCount++
		default:
			// Channel is full, drop the event
			m.logger.Warn("Channel is full, dropping event",
				zap.String("topic", topic),
				zap.String("client_id", clientID),
			)
		}
	}

	m.logger.Debug("Event broadcast completed",
		zap.String("topic", topic),
		zap.Int("total_clients", clientCount),
		zap.Int("successful_deliveries", successCount),
		zap.Int("dropped_events", clientCount-successCount),
	)
}

func (m *Manager[T]) GetClient(topic string, clientID string) chan T {
	if clients, exists := m.clients[topic]; exists {
		if clientChan, exists := clients[clientID]; exists {
			return clientChan
		}
	}

	return nil
}

// Cleanup removes all clients from a topic
func (m *Manager[T]) Cleanup(topic string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if clients, exists := m.clients[topic]; exists {
		for _, clientChan := range clients {
			close(clientChan)
		}
		delete(m.clients, topic)
	}
}

// generateClientID generates a unique client ID
func generateClientID() string {
	// Use nanosecond timestamp + random suffix for maximum uniqueness
	timestamp := time.Now().UTC().UnixNano()
	randomSuffix := fmt.Sprintf("%06d", time.Now().UTC().UnixNano()%1000000)
	return fmt.Sprintf("client_%d_%s", timestamp, randomSuffix)
}
