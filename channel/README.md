# Channel Manager

A generic, thread-safe channel manager for handling real-time event broadcasting in Go applications.

## Features

- **Generic**: Works with any event type using Go generics
- **Thread-safe**: Built-in mutex handling for concurrent operations
- **Automatic cleanup**: Removes empty topics when all clients disconnect
- **Configurable**: Customizable buffer sizes for channels
- **Logging**: Built-in logging support for debugging and monitoring
- **Simple API**: Easy to use with minimal boilerplate

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "time"
    "github.com/TixiaOTA/tixia-backend-master/pkg/channel"
    "github.com/TixiaOTA/tixia-backend-master/pkg/logger"
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

### API Reference

#### `NewManager[T](bufferSize int, logger *logger.Logger) *Manager[T]`
Creates a new channel manager with the specified buffer size for channels and logger instance.

#### `RegisterClient(topic string, clientID string) (string, chan T)`
Registers a new client for a topic and returns the client ID and channel. If clientID is empty, a unique ID will be generated.

#### `UnregisterClient(topic, clientID string)`
Removes a client from a topic and closes its channel. Automatically cleans up empty topics.

#### `Broadcast(topic string, event T)`
Sends an event to all clients subscribed to a topic. Events are dropped if a client's channel is full.

#### `GetClient(topic string, clientID string) chan T`
Returns the channel for a specific client in a topic, or nil if not found.

#### `Cleanup(topic string)`
Removes all clients from a topic and closes their channels.

## Thread Safety

All operations are thread-safe and can be called from multiple goroutines concurrently.

## Error Handling

The manager handles common scenarios gracefully:
- Full channels: Events are dropped when a client's channel is full (logged as warning)
- Disconnected clients: Channels are automatically closed and cleaned up
- Non-existent topics: Operations on non-existent topics are safely ignored
- Automatic client ID generation: If no client ID is provided, a timestamp-based ID is generated

## Logging

The manager logs important events:
- Client registration with topic and client ID
- Warnings when channels are full and events are dropped
