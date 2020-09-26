package util

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func GetOriginHost(c *gin.Context) string {
	origin := c.Request.Header.Get("Origin")
	refer := c.Request.Header.Get("Referer")
	if origin != "" {
		index := strings.Index(origin, "//")
		return origin[index+2:]
	} else {
		index := strings.Index(refer, "//")
		subRefer := refer[index+2:]
		index2 := strings.Index(subRefer, "/")
		return subRefer[0:index2]
	}

}
