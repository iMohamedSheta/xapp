//go:build !dev

package hooks

import (
	"context"

	"github.com/gin-gonic/gin"
)

func attachDebugProps(ctx context.Context, c *gin.Context, reqID string) {}
