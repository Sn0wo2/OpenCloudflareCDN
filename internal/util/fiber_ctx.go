package util

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func FiberContextString(ctx *fiber.Ctx) string {
	var sb strings.Builder

	ips := ctx.IPs()
	if len(ips) == 0 {
		ips = []string{ctx.IP()}
	}

	sb.WriteString(strings.Join(ips, ", "))

	sb.WriteString(" -> ")
	sb.WriteString(ctx.Method())

	sb.WriteString(" ")

	if ctx.Response().StatusCode() != 0 {
		statusCode := ctx.Response().StatusCode()
		sb.WriteString(strconv.Itoa(statusCode))
		sb.WriteString(" ")
		sb.WriteString(http.StatusText(statusCode))
		sb.WriteString(" ")
	}

	sb.WriteString(BytesToString(ctx.Request().RequestURI()))

	var headers []string

	ctx.Request().Header.All()(func(key, value []byte) bool {
		v := BytesToString(value)
		if len(v) > 20 {
			v = v[:12] + "..."
		}

		headers = append(headers, fmt.Sprintf("%s:%s", BytesToString(key), v))

		return true
	})

	if len(headers) > 0 {
		sb.WriteString(" { ")
		sb.WriteString(strings.Join(headers, ", "))
		sb.WriteString(" }")
	}

	return sb.String()
}
