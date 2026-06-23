package middleware

import (
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/imohamedsheta/xapp/app/shared/enums"
	"github.com/imohamedsheta/xapp/app/shared/hooks"
	"github.com/imohamedsheta/xapp/app/shared/utils"
	"github.com/imohamedsheta/xapp/app/x"
	"go.uber.org/zap"
)

type debugResponseWriter struct {
	gin.ResponseWriter
	c             *gin.Context
	reqID         string
	headerWritten bool
}

func (w *debugResponseWriter) WriteHeader(statusCode int) {
	w.attachHeader()
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *debugResponseWriter) Write(b []byte) (int, error) {
	w.attachHeader()
	return w.ResponseWriter.Write(b)
}

func (w *debugResponseWriter) WriteString(s string) (int, error) {
	w.attachHeader()
	return w.ResponseWriter.WriteString(s)
}

func (w *debugResponseWriter) attachHeader() {
	if w.headerWritten {
		return
	}
	w.headerWritten = true
	hooks.AttachDebugHeader(w.c, w.reqID)
}

// Logger is a Gin middleware for logging HTTP requests using zap
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Generate or get request ID for tracing
		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = uuid.NewString()
			c.Request.Header.Set("X-Request-ID", reqID)
		}
		c.Writer.Header().Set("X-Request-ID", reqID)
		c.Set(enums.ContextKeyRequestId.String(), reqID)

		if x.Config().GetBool("app.debug", true) {
			c.Set(enums.ContextKeyRequestStartTime.String(), time.Now().UnixMilli())
			c.Set("request_method", c.Request.Method)
			c.Set("request_path", c.FullPath()) // or c.Request.URL.Path
			c.Set("request_query", c.Request.URL.RawQuery)
			c.Set("request_client_ip", c.ClientIP())
			c.Set("request_user_agent", c.Request.UserAgent())
			c.Set("request_referer", c.Request.Referer())
			c.Set("request_content_length", c.Request.ContentLength)
			c.Set("request_host", c.Request.Host)
			c.Set("request_protocol", c.Request.Proto)
			var memStart runtime.MemStats
			runtime.ReadMemStats(&memStart)
			c.Set("memory_start", memStart)
		}

		path := c.Request.URL.Path
		utils.PrintErr("Path: " + path + " - request ID: " + reqID)

		// Wrap response writer to capture and inject headers before response is flushed
		dw := &debugResponseWriter{
			ResponseWriter: c.Writer,
			c:              c,
			reqID:          reqID,
		}
		c.Writer = dw

		// Process request
		c.Next()

		// Just in case no body was written (e.g. 204 No Content, 302 Redirect), trigger header writing now
		dw.attachHeader()

		// request termination (after the request has been processed)
		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		// path := c.Request.URL.Path
		userAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		contentLength := c.Request.ContentLength
		host := c.Request.Host
		protocol := c.Request.Proto
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()
		errorCount := len(c.Errors)

		log := x.Log().Channel("request_log").With(
			zap.Int("status", status),
			zap.String("latency", latency.String()),
			zap.String("client_ip", clientIP),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("user_agent", userAgent),
			zap.String("referer", referer),
			zap.Int64("content_length", contentLength),
			zap.String("host", host),
			zap.String("protocol", protocol),
			zap.String("request_id", reqID),
			zap.Int("error_count", errorCount),
		)

		if errorMessage != "" {
			log.Error(errorMessage)
		} else {
			log.Info("Handled request")
		}
	}
}
