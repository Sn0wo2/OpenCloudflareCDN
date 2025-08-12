package handler

import (
	"fmt"

	"github.com/Sn0wo2/OpenCloudflareCDN/config"
	"github.com/Sn0wo2/OpenCloudflareCDN/internal/util"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/valyala/fasthttp"
)

func Proxy() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tokenStr := ctx.Cookies("cfv_c")
		if tokenStr != "" {
			_, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				return util.StringToBytes(config.Instance.JWTSecret), nil
			})
			if err == nil {
				return proxyRequest(ctx)
			}
		}

		return ctx.Next()
	}
}

func proxyRequest(ctx *fiber.Ctx) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	ctx.Request().CopyTo(req)

	req.SetRequestURI(config.Instance.OriginalServer + ctx.OriginalURL())

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	client := fasthttp.Client{}

	if err := client.Do(req, resp); err != nil {
		return ctx.Status(fiber.StatusBadGateway).SendString("Bad Gateway: " + err.Error())
	}

	ctx.Status(resp.StatusCode())

	for k, v := range resp.Header.All() {
		ctx.Response().Header.SetBytesKV(k, v)
	}

	ctx.Response().SetBody(resp.Body())

	return nil
}
