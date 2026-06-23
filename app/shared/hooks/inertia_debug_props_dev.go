//go:build !dev

package hooks

import (
	"github.com/gin-gonic/gin"
	"github.com/imohamedsheta/xapp/app/shared/utils"
)

// AttachDebugHeader is a no-op in non-dev builds.
func AttachDebugHeader(c *gin.Context, reqID string) {
	utils.PrintErr("AttachDebugHeader: called non-dev no-op stub for request: " + reqID)
}

// AttachDebugProps is a no-op in non-dev builds.
func AttachDebugProps(c *gin.Context, reqID string) {}
