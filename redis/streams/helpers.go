package streams

import (
	"context"
	"fmt"
	"time"

	"github.com/TixiaOTA/tixia-pkg/logger"

	"go.uber.org/zap"
)

func GetEventMetadata(event Event, key string) (string, bool) {
	metadata, ok := event.Metadata[key].(string)
	if !ok {
		return "", false
	}
	return metadata, true
}

func GetEventData(event Event, key string) (map[string]any, bool) {
	data, ok := event.Data[key].(map[string]any)
	if !ok {
		return nil, false
	}
	return data, true
}

func GenerateEventID() string {
	return fmt.Sprintf("evt_%d", time.Now().UTC().UnixNano())
}

func StartEventConsumer(ctx context.Context, streamManager *StreamManager, streamPrefix, topic string, log *logger.Logger, timeout time.Duration) {
	streamName := fmt.Sprintf("%s:%s", streamPrefix, topic)
	go func() {
		if err := streamManager.StartConsumer(ctx, streamName, timeout); err != nil {
			log.Error(fmt.Sprintf("%s consumer failed", topic), zap.Error(err))
		}
	}()
}

// StartEventConsumerBatch starts a single goroutine consuming multiple topics in one loop
func StartEventConsumerBatch(ctx context.Context, streamManager *StreamManager, streamPrefix string, topics []string, log *logger.Logger, timeout time.Duration) {
	streamNames := make([]string, 0, len(topics))
	for _, t := range topics {
		streamNames = append(streamNames, fmt.Sprintf("%s:%s", streamPrefix, t))
	}
	go func() {
		if err := streamManager.StartConsumerForStreams(ctx, streamNames, timeout); err != nil {
			log.Error("batch consumer failed", zap.Error(err))
		}
	}()
}
