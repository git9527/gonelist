package util

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func GetOriginHost(c *gin.Context) string {
	origin := c.Request.Header.Get("Origin")
	index := strings.Index(origin, "//")
	return origin[index+2:]
}
