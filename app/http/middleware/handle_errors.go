package middleware

import (
	"errors"
	"fmt"

	"net"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iMohamedSheta/xerr"
	"github.com/imohamedsheta/xapp/app/shared/enums"
	"github.com/imohamedsheta/xapp/app/shared/utils"
	"github.com/imohamedsheta/xapp/app/x"
	"github.com/imohamedsheta/xapp/pkg/inertia"
	"go.uber.org/zap"
)

type HandlerFuncWithError func(*gin.Context) error

func HandleErrors(handler HandlerFuncWithError) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := handler(c); err != nil {
			ErrorHandler(c, err, false)
			return
		}
	}
}

func ErrorHandler(c *gin.Context, err error, isInertiaReq bool) {
	globalErrorHandler(c, err)

	var xe *xerr.XErr
	if ok := errors.As(err, &xe); ok && xe != nil {
		switch {
		case xe.IsType(enums.XErrBadRequestError):
			badRequestErrorHandler(c, xe, isInertiaReq)
		case xe.IsType(enums.XErrValidationError):
			validationErrorHandler(c, xe, isInertiaReq)
		case xe.IsType(enums.XErrBadRequestBindingError):
			badRequestBindingErrorHandler(c, xe, isInertiaReq)
		case xe.IsType(enums.XErrUnAuthorizedError):
			unAuthorizedErrorHandler(c, xe, isInertiaReq)
		case xe.IsType(enums.XErrForbiddenError):
			forbiddenErrorHandler(c, xe, isInertiaReq)
		case xe.IsType(enums.XErrNotFoundError):
			notFoundErrorHandler(c, xe, isInertiaReq)
		case xe.IsType(enums.XErrServerError):
			serverErrorHandler(c, xe, isInertiaReq)
		default:
			x.Logger().Error("Unhandled error type xerr ", zap.Error(xe))
			unknownErrorHandler(c, xe, isInertiaReq)
		}
	} else {
		unknownErrorHandler(c, err, isInertiaReq)
	}

	c.Abort()
}

func validationErrorHandler(c *gin.Context, err *xerr.XErr, isInertiaReq bool) {
	errorResponse(c, err, isInertiaReq)
}

func badRequestErrorHandler(c *gin.Context, err *xerr.XErr, isInertiaReq bool) {
	errorResponse(c, err, isInertiaReq)
}

func badRequestBindingErrorHandler(c *gin.Context, err *xerr.XErr, isInertiaReq bool) {
	errorResponse(c, err, isInertiaReq)
}

func notFoundErrorHandler(c *gin.Context, err *xerr.XErr, isInertiaReq bool) {
	errorResponse(c, err, isInertiaReq)
}

func unAuthorizedErrorHandler(c *gin.Context, err *xerr.XErr, isInertiaReq bool) {
	if isInertiaReq {
		x.Inertia().Redirect(c.Writer, c.Request, "/login")
	}
	errorResponse(c, err, false)
}

func forbiddenErrorHandler(c *gin.Context, err *xerr.XErr, isInertiaReq bool) {
	errorResponse(c, err, isInertiaReq)
}

func serverErrorHandler(c *gin.Context, err *xerr.XErr, isInertiaReq bool) {
	errorResponse(c, err, isInertiaReq)
}

func unknownErrorHandler(c *gin.Context, err error, isInertiaReq bool) {
	serverErrorHandler(c, xerr.New("unknown error", enums.XErrServerError, err), isInertiaReq)
}

func errorResponse(c *gin.Context, err *xerr.XErr, isInertiaRequest bool) {
	code := enums.GetErrorCode(err.Type)
	errMsg := err.PublicMessage
	if errMsg == "" {
		errMsg = code.MessageAr()
	}

	if isInertiaRequestOrBrowserVisit(c) || isInertiaRequest {
		i := x.Inertia()
		// For GET requests, render an error page
		if c.Request.Method == "GET" {
			renderErr := i.Render(c, "Errors/Error", inertia.Props{
				"message":           errMsg,
				"error_code":        fmt.Sprintf("%d %s", code.StatusCode(), code.String()),
				"error_status_code": code.StatusCode(),
			})

			if renderErr != nil {
				utils.Dump("There is render error (error)")
				// API/JSON response for non-Inertia requests
				data := map[string]any{
					"message":    errMsg,
					"error_code": code.String(),
				}

				if err.IsType(enums.XErrValidationError) {
					data["errors"] = err.Details
				}

				c.JSON(code.StatusCode(), data)
			}
			return
		}

		// For non-GET requests (POST, PUT, DELETE), flash alert and redirect back
		flashProps := inertia.Props{
			"toast": []map[string]any{utils.ToastError("خطأ", errMsg)},
		}
		if flashErr := x.Flash().Flash(c, flashProps); flashErr != nil {
			x.Logger().Error(flashErr.Error())
		}
		x.Inertia().Back(c, code.StatusCode())
		return
	}

	// API/JSON response for non-Inertia requests
	data := map[string]any{
		"message":    errMsg,
		"error_code": code.String(),
	}

	if err.IsType(enums.XErrValidationError) {
		data["errors"] = err.Details
	}

	c.JSON(code.StatusCode(), data)
}

func globalErrorHandler(c *gin.Context, err error) {
	appRequestErrorLogger(c, err)
}

func appRequestErrorLogger(c *gin.Context, err error) {
	// Get stack trace
	stackBuf := make([]byte, 1024*8) // 8KB buffer
	stackSize := runtime.Stack(stackBuf, false)
	stack := string(stackBuf[:stackSize])

	// Get caller information
	_, file, line, ok := runtime.Caller(2)
	var caller string
	if ok {
		// Extract just the filename
		parts := strings.Split(file, "/")
		if len(parts) > 0 {
			caller = parts[len(parts)-1]
		}
	}

	requestID := c.GetString("request_id")
	if requestID == "" {
		requestID = c.GetHeader("X-Request-ID")
	}

	userID := c.GetString(string(enums.ContextKeyAuthId))
	if userID == "" {
		userID = "anonymous"
	}

	log := x.Logger()

	// Log comprehensive error information
	errStr := "unknown error"
	if err != nil {
		errStr = err.Error()
	}

	log.Error("Request error occurred",
		// Error
		zap.String("error", errStr),
		// zap.String("error_type", getErrorType(err)),

		// Request context
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("raw_query", c.Request.URL.RawQuery),
		zap.String("user_agent", c.GetHeader("User-Agent")),
		zap.String("client_ip", c.ClientIP()),
		zap.String("request_id", requestID),
		zap.String("user_id", userID),

		// Request headers
		zap.Any("headers", sanitizeHeaders(c.Request.Header)),

		// Timing
		zap.String("timestamp", time.Now().Format("2006-01-02 15:04:05")),

		// Code location
		zap.String("file", caller),
		zap.Int("line", line),

		// Stack trace (you might want to make this conditional based on log level)
		zap.String("stack_trace", stack),

		// Additional context if available
		zap.Any("request_body_size", c.Request.ContentLength),
		zap.String("content_type", c.GetHeader("Content-Type")),
		zap.String("accept", c.GetHeader("Accept")),
		zap.String("referer", c.GetHeader("Referer")),
	)

	// Also log request parameters if they exist
	if len(c.Params) > 0 {
		params := make(map[string]string)
		for _, param := range c.Params {
			params[param.Key] = param.Value
		}
		log.Info("Request parameters",
			zap.Any("params", params),
			zap.String("request_id", requestID),
		)
	}

	// Log query parameters if they exist
	if len(c.Request.URL.Query()) > 0 {
		log.Info("Query parameters",
			zap.Any("query_params", sanitizeQueryParams(c.Request.URL.Query())),
			zap.String("request_id", requestID),
		)
	}
}

func RequestErrorLogger(w http.ResponseWriter, r *http.Request, err error) {
	// Get stack trace
	stackBuf := make([]byte, 8*1024)
	stackSize := runtime.Stack(stackBuf, false)
	stack := string(stackBuf[:stackSize])

	// Get caller info
	_, file, line, ok := runtime.Caller(2)
	var caller string
	if ok {
		parts := strings.Split(file, "/")
		if len(parts) > 0 {
			caller = parts[len(parts)-1]
		}
	}

	// Extract request ID from headers (no Gin context)
	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = "unknown"
	}

	// User ID won't be available without Gin context (or custom headers)
	userID := "anonymous"

	log := x.Logger()

	log.Error("Request error occurred",
		zap.String("error", err.Error()),
		// zap.String("error_type", getErrorType(err)),

		// Request info
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("raw_query", r.URL.RawQuery),
		zap.String("user_agent", r.Header.Get("User-Agent")),
		zap.String("client_ip", clientIPFromRequest(r)),

		zap.String("request_id", requestID),
		zap.String("user_id", userID),

		zap.Any("headers", sanitizeHeaders(r.Header)),

		zap.String("timestamp", time.Now().Format("2006-01-02 15:04:05")),

		zap.String("file", caller),
		zap.Int("line", line),

		zap.String("stack_trace", stack),

		zap.Int64("request_body_size", r.ContentLength),
		zap.String("content_type", r.Header.Get("Content-Type")),
		zap.String("accept", r.Header.Get("Accept")),
		zap.String("referer", r.Header.Get("Referer")),
	)
}

func isInertiaRequestOrBrowserVisit(c *gin.Context) bool {
	acceptHeader := c.GetHeader("Accept")
	isBrowserVisit := strings.Contains(acceptHeader, "text/html")
	return inertia.IsInertiaRequest(c.Request) || isBrowserVisit
}

func clientIPFromRequest(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// Helper function to sanitize headers
func sanitizeHeaders(headers map[string][]string) map[string][]string {
	sanitized := make(map[string][]string)
	sensitiveHeaders := map[string]bool{
		"authorization": true,
		"cookie":        true,
		"x-api-key":     true,
		"x-auth-token":  true,
	}

	for key, values := range headers {
		lowerKey := strings.ToLower(key)
		if sensitiveHeaders[lowerKey] {
			sanitized[key] = []string{"[REDACTED]"}
		} else {
			sanitized[key] = values
		}
	}
	return sanitized
}

// Helper function to sanitize query parameters
func sanitizeQueryParams(params map[string][]string) map[string][]string {
	sanitized := make(map[string][]string)
	sensitiveParams := map[string]bool{
		"password": true,
		"api_key":  true,
		"token":    true,
		"secret":   true,
	}

	for key, values := range params {
		lowerKey := strings.ToLower(key)
		if sensitiveParams[lowerKey] {
			sanitized[key] = []string{"[REDACTED]"}
		} else {
			sanitized[key] = values
		}
	}
	return sanitized
}
