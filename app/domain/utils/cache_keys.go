package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetAuthSessionSettingsCacheKey(c *gin.Context) (string, error) {
	sessionId, err := AuthSessionId(c)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("session:settings:%s", sessionId), nil
}
