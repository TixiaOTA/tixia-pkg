package redisutil

import (
	"encoding/json"
	"fmt"
)

// redisutil/trace_context.go
func NormalizeTraceMap(raw interface{}) map[string]string {
	traceMap := make(map[string]string)
	switch t := raw.(type) {
	case map[string]interface{}:
		for k, v := range t {
			traceMap[k] = fmt.Sprintf("%v", v)
		}
	case string:
		_ = json.Unmarshal([]byte(t), &traceMap)
	}
	return traceMap
}
