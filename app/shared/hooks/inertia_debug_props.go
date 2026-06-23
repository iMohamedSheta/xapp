//go:build dev

package hooks

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imohamedsheta/xapp/app/shared/enums"
	"github.com/imohamedsheta/xapp/app/x"
)

type requestHistoryStore struct {
	mu       sync.RWMutex
	requests []map[string]any
}

var historyStore = &requestHistoryStore{
	requests: make([]map[string]any, 0, 15),
}

func addRequestToHistory(data map[string]any) {
	historyStore.mu.Lock()
	defer historyStore.mu.Unlock()

	reqID := ""
	if req, ok := data["request"].(map[string]any); ok {
		if id, ok := req["id"].(string); ok {
			reqID = id
		}
	}

	if reqID != "" {
		for i, r := range historyStore.requests {
			if req, ok := r["request"].(map[string]any); ok {
				if id, ok := req["id"].(string); ok && id == reqID {
					historyStore.requests[i] = data
					return
				}
			}
		}
	}

	historyStore.requests = append([]map[string]any{data}, historyStore.requests...)
	if len(historyStore.requests) > 15 {
		historyStore.requests = historyStore.requests[:15]
	}
}

func getRequestHistory() []map[string]any {
	historyStore.mu.RLock()
	defer historyStore.mu.RUnlock()

	res := make([]map[string]any, len(historyStore.requests))
	copy(res, historyStore.requests)
	return res
}

// GetDebugData compiles the query, request, and route debug information.
func GetDebugData(c *gin.Context, reqID string) map[string]any {
	var requestMemoryUsed uint64
	var memEnd runtime.MemStats
	runtime.ReadMemStats(&memEnd)

	memStartObj := c.Value("memory_start")
	if memStartObj != nil {
		if memStart, ok := memStartObj.(runtime.MemStats); ok {
			requestMemoryUsed = memEnd.Alloc - memStart.Alloc
		}
	}

	action, handler, short, middlewares := extractRouteInfo(c)

	data := map[string]any{
		"queries": GetQueries(reqID),
		"request": map[string]any{
			"id":                       reqID,
			"method":                   c.Value("request_method"),
			"path":                     c.Value("request_path"),
			"query":                    c.Value("request_query"),
			"client_ip":                c.Value("request_client_ip"),
			"user_agent":               c.Value("request_user_agent"),
			"referer":                  c.Value("request_referer"),
			"content_length":           c.Value("request_content_length"),
			"host":                     c.Value("request_host"),
			"protocol":                 c.Value("request_protocol"),
			"action_duration":          calculateDuration(c),
			"status":                   c.Writer.Status(),
			"memory_usage":             fmt.Sprintf("%.2f MB", float64(memEnd.Alloc/1024/1024)),
			"memory_usage_for_request": fmt.Sprintf("%.2f KB", float64(requestMemoryUsed/1024)),
		},
		"route": map[string]any{
			"uri":        c.FullPath(),
			"handler":    handler,
			"method":     c.Request.Method,
			"action":     action,
			"short":      short,
			"middleware": middlewares,
		},
	}

	addRequestToHistory(data)

	// Create a shallow copy of the map to avoid JSON serialization cycles.
	// This ensures the "history" slice is only attached to the returned copy,
	// and doesn't pollute the entries in historyStore.
	response := make(map[string]any, len(data)+1)
	for k, v := range data {
		response[k] = v
	}
	response["history"] = getRequestHistory()

	return response
}

// AttachDebugHeader is the public entry point called by the logger middleware.
// It collects debug info and writes it as an X-Debug-Data response header,
// making it available for all HTTP methods (GET, POST, PATCH, DELETE, etc.).
func AttachDebugHeader(c *gin.Context, reqID string) {
	debugData := GetDebugData(c, reqID)

	// Attach as response header so it works for all HTTP methods (GET, POST, PATCH, DELETE…)
	if b, err := json.Marshal(debugData); err == nil {
		c.Header("X-Debug-Data", string(b))
	}

	ClearQueries(reqID)
}

// AttachDebugProps is called before rendering Inertia pages to populate page.props.debug
func AttachDebugProps(c *gin.Context, reqID string) {
	debugData := GetDebugData(c, reqID)
	// We use direct package call here to share props.
	// Since we import "github.com/imohamedsheta/xapp/app/x", we must add it to the imports
	// Note: We don't ClearQueries here because the logger middleware will run after this and clear it.
	// That way queries are captured in both shared props and the final response header.
	x.Inertia().ShareProp("debug", debugData)
}

func extractRouteInfo(c *gin.Context) (action string, handler string, short string, middlewares []string) {
	all := c.HandlerNames()
	middlewares = []string{}

	// use stored values if present (set in wI / InertiaMiddlewareWithErrorHandler)
	if act, ok := c.Get("route_action_full"); ok {
		action = shortActionName(act.(string))
	}
	if sh, ok := c.Get("route_action_short"); ok {
		short = shortActionName(sh.(string))
	}

	// fallback if nothing was stored
	if action == "" && len(all) > 0 {
		action = shortActionName(all[len(all)-1])
		short = shortActionName(action)
	}

	// collect middleware names
	seen := make(map[string]bool)
	for i, fullname := range all {
		// skip action
		if fullname == action || i == len(all)-1 {
			continue
		}

		name := fullname
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true
		middlewares = append(middlewares, name)
	}

	handler = strings.TrimSuffix(action, "."+short)
	handler = strings.TrimPrefix(handler, "imohamedsheta/xapp/")
	return action, handler, short, middlewares
}

func shortActionName(full string) string {
	// strip Go's method wrapper suffixes
	full = strings.TrimSuffix(full, "-fm")
	full = strings.TrimSuffix(full, "·fm")

	full = strings.TrimPrefix(full, "imohamedsheta/xapp/")

	// cleanup receivers: (*EventHandler) → EventHandler
	full = strings.ReplaceAll(full, "(*", "")
	full = strings.ReplaceAll(full, ")", "")
	return full
}

func calculateDuration(c *gin.Context) string {
	startTime := c.Value(enums.ContextKeyRequestStartTime.String()).(int64)
	return time.Since(time.UnixMilli(startTime)).String()
}
