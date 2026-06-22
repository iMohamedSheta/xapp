//go:build dev

package hooks

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imohamedsheta/xapp/app/domain/enums"
	"github.com/imohamedsheta/xapp/app/x"
)

func attachDebugProps(ctx context.Context, c *gin.Context, reqID string) {
	memStart := ctx.Value("memory_start").(runtime.MemStats)
	mem := &runtime.MemStats{}
	runtime.ReadMemStats(mem)

	memEnd := &runtime.MemStats{}
	runtime.ReadMemStats(memEnd)
	requestMemoryUsed := memEnd.Alloc - memStart.Alloc

	action, handler, short, middlewares := extractRouteInfo(c)

	x.Inertia().ShareProp("debug", map[string]any{
		"queries": GetQueries(reqID),
		"request": map[string]any{
			"id":                       reqID,
			"method":                   ctx.Value("request_method"),
			"path":                     ctx.Value("request_path"),
			"query":                    ctx.Value("request_query"),
			"client_ip":                ctx.Value("request_client_ip"),
			"user_agent":               ctx.Value("request_user_agent"),
			"referer":                  ctx.Value("request_referer"),
			"content_length":           ctx.Value("request_content_length"),
			"host":                     ctx.Value("request_host"),
			"protocol":                 ctx.Value("request_protocol"),
			"action_duration":          calculateDuration(ctx),
			"status":                   c.Writer.Status(),
			"memory_usage":             fmt.Sprintf("%.2f MB", float64(mem.Alloc/1024/1024)),
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
	})
	defer func() {
		ClearQueries(reqID)
	}()
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
	// strip Go’s method wrapper suffixes
	full = strings.TrimSuffix(full, "-fm")
	full = strings.TrimSuffix(full, "·fm")

	full = strings.TrimPrefix(full, "imohamedsheta/xapp/")

	// cleanup receivers: (*EventHandler) → EventHandler
	full = strings.ReplaceAll(full, "(*", "")
	full = strings.ReplaceAll(full, ")", "")
	return full
}

func calculateDuration(ctx context.Context) string {
	startTime := ctx.Value(enums.ContextKeyRequestStartTime.String()).(int64)
	return time.Since(time.UnixMilli(startTime)).String()
}
