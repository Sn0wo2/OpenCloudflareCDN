package util

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GinContextString(ctx *gin.Context) string {
	var sb strings.Builder

	ips := ctx.ClientIP()

	sb.WriteString(ips)

	sb.WriteString(" -> ")
	sb.WriteString(ctx.Request.Method)

	sb.WriteString(" ")

	if ctx.Writer.Status() != 0 {
		statusCode := ctx.Writer.Status()
		sb.WriteString(strconv.Itoa(statusCode))
		sb.WriteString(" ")
		sb.WriteString(http.StatusText(statusCode))
		sb.WriteString(" ")
	}

	sb.WriteString(ctx.Request.RequestURI)

	var headers []string

	for key, values := range ctx.Request.Header {
		for _, value := range values {
			v := value
			if len(v) > 20 {
				v = v[:12] + "..."
			}
			headers = append(headers, fmt.Sprintf("%s:%s", key, v))
		}
	}

	if len(headers) > 0 {
		sb.WriteString(" { ")
		sb.WriteString(strings.Join(headers, ", "))
		sb.WriteString(" }")
	}

	return sb.String()
}
