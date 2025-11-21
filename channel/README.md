# Channel Manager Package

A generic, thread-safe channel manager for handling real-time event broadcasting in Go applications.

## Overview

The Channel Manager package enables you to:

- Broadcast events to multiple clients in real-time
- Manage client subscriptions to topics
- Handle concurrent operations safely with built-in thread safety
- Automatically clean up empty topics when clients disconnect
- Use Go generics to work with any event type

## Installation

```bash
go get github.com/TixiaOTA/tixia-pkg/channel
```

## Configuration

### Configuration

The channel manager is initialized with a buffer size and logger instance:

```go
import (
    "github.com/TixiaOTA/tixia-pkg/channel"
    "github.com/TixiaOTA/tixia-pkg/logger"
)

logger := logger.NewLogger()
manager := channel.NewManager[YourEventType](100, logger)
```

### Configuration Fields

- **bufferSize** (int, required): Buffer size for client channels. If <= 0, defaults to 100
- **logger** (\*logger.Logger, required): Logger instance for logging operations

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "github.com/TixiaOTA/tixia-pkg/channel"
    "github.com/TixiaOTA/tixia-pkg/logger"
)

type FlightEvent struct {
    SearchID string `json:"search_id"`
    Status   string `json:"status"`
    Data     []Flight `json:"data"`
}

func main() {
    // Initialize logger
    logger := logger.NewLogger()

    // Create a channel manager for FlightEvent type with buffer size 100
    manager := channel.NewManager[FlightEvent](100, logger)

    // Register a client for a specific search
    clientID, clientChan := manager.RegisterClient("search_123", "")
    defer manager.UnregisterClient("search_123", clientID)

    // Listen for events
    go func() {
        for event := range clientChan {
            fmt.Printf("Received flight event: %+v\n", event)
        }
    }()

    // Broadcast an event
    event := FlightEvent{
        SearchID: "search_123",
        Status:   "completed",
        Data:     []Flight{},
    }
    manager.Broadcast("search_123", event)
}
```

### Multiple Clients Example

```go
// Register multiple clients for the same topic
client1ID, chan1 := manager.RegisterClient("search_123", "")
client2ID, chan2 := manager.RegisterClient("search_123", "")

// Broadcast to all clients
manager.Broadcast("search_123", event)

// Both clients will receive the event
```

## API Reference

### Public Functions

#### NewManager

Creates a new channel manager with the specified buffer size and logger.

```go
func NewManager[T any](bufferSize int, logger *logger.Logger) *Manager[T]
```

**Parameters:**

- `bufferSize` (int, required): Buffer size for client channels. Defaults to 100 if <= 0
- `logger` (\*logger.Logger, required): Logger instance for logging operations

**Return Values:**

- `*Manager[T]`: A new channel manager instance

**Special Notes:** The manager uses Go generics, so you must specify the event type when creating it.

### Public Types

#### Manager

The main manager struct that handles client subscriptions and event broadcasting.

```go
type Manager[T any] struct {
    // Private fields: clients, mutex, bufferSize, logger
}
```

#### Methods on Manager

##### RegisterClient

Registers a new client for a topic and returns the client ID and channel.

```go
func (m *Manager[T]) RegisterClient(topic string, clientID string) (string, chan T)
```

**Parameters:**

- `topic` (string, required): Topic name to subscribe to
- `clientID` (string, optional): Client ID. If empty, a unique ID will be generated

**Return Values:**

- `string`: Client ID (generated if not provided)
- `chan T`: Channel for receiving events

**Special Notes:** If a client ID collision is detected, a new unique ID will be automatically generated.

##### UnregisterClient

Removes a client from a topic and closes its channel.

```go
func (m *Manager[T]) UnregisterClient(topic, clientID string)
```

**Parameters:**

- `topic` (string, required): Topic name
- `clientID` (string, required): Client ID to remove

**Special Notes:** Automatically cleans up empty topics when all clients disconnect.

##### Broadcast

Sends an event to all clients subscribed to a topic.

```go
func (m *Manager[T]) Broadcast(topic string, event T)
```

**Parameters:**

- `topic` (string, required): Topic name to broadcast to
- `event` (T, required): Event to broadcast

**Special Notes:** Events are dropped if a client's channel is full (logged as warning).

##### GetClient

Returns the channel for a specific client in a topic.

```go
func (m *Manager[T]) GetClient(topic string, clientID string) chan T
```

**Parameters:**

- `topic` (string, required): Topic name
- `clientID` (string, required): Client ID

**Return Values:**

- `chan T`: Client channel, or nil if not found

##### Cleanup

Removes all clients from a topic and closes their channels.

```go
func (m *Manager[T]) Cleanup(topic string)
```

**Parameters:**

- `topic` (string, required): Topic name to clean up

**Special Notes:** Useful for shutting down a topic completely.

## Error Handling

The manager handles common scenarios gracefully:

- **Full channels**: Events are dropped when a client's channel is full (logged as warning)
- **Disconnected clients**: Channels are automatically closed and cleaned up
- **Non-existent topics**: Operations on non-existent topics are safely ignored
- **Client ID collisions**: If a client ID collision is detected, a new unique ID is automatically generated

## Data Models

### Manager

The main manager struct:

```go
type Manager[T any] struct {
    clients    map[string]map[string]chan T // topic -> clientID -> channel
    mutex      sync.RWMutex                  // Thread safety
    bufferSize int                           // Channel buffer size
    logger     *logger.Logger                // Logger instance
}
```

The manager uses Go generics (`[T any]`) to work with any event type. The internal structure maps topics to clients, where each client has a unique ID and a channel for receiving events.

## Examples

### Complete Example

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/TixiaOTA/tixia-pkg/channel"
    "github.com/TixiaOTA/tixia-pkg/logger"
)

type SearchEvent struct {
    SearchID string    `json:"search_id"`
    Status   string    `json:"status"`
    Progress int       `json:"progress"`
    Timestamp time.Time `json:"timestamp"`
}

func main() {
    // Initialize logger
    logger := logger.NewLogger()

    // Create manager with buffer size 100
    manager := channel.NewManager[SearchEvent](100, logger)

    // Register multiple clients for the same topic
    client1ID, chan1 := manager.RegisterClient("search_123", "")
    client2ID, chan2 := manager.RegisterClient("search_123", "")

    defer func() {
        manager.UnregisterClient("search_123", client1ID)
        manager.UnregisterClient("search_123", client2ID)
    }()

    // Start listening goroutines
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    go listenForEvents(ctx, "Client 1", chan1)
    go listenForEvents(ctx, "Client 2", chan2)

    // Simulate broadcasting events
    for i := 0; i < 5; i++ {
        event := SearchEvent{
            SearchID:  "search_123",
            Status:    "processing",
            Progress:  i * 20,
            Timestamp: time.Now(),
        }

        manager.Broadcast("search_123", event)
        time.Sleep(500 * time.Millisecond)
    }

    // Final event
    manager.Broadcast("search_123", SearchEvent{
        SearchID:  "search_123",
        Status:    "completed",
        Progress:  100,
        Timestamp: time.Now(),
    })

    time.Sleep(1 * time.Second)
}

func listenForEvents(ctx context.Context, clientName string, ch <-chan SearchEvent) {
    for {
        select {
        case event := <-ch:
            fmt.Printf("[%s] Received: %+v\n", clientName, event)
        case <-ctx.Done():
            return
        }
    }
}
```

## License

This package is part of the TixiaOTA project.
