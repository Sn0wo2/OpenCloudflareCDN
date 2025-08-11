package handler

import (
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/Sn0wo2/OpenCloudflareCDN/config"
	"github.com/Sn0wo2/OpenCloudflareCDN/internal/util"
	"github.com/Sn0wo2/OpenCloudflareCDN/response"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tidwall/gjson"
)

func APIVerify() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		resp, err := http.PostForm("https://challenges.cloudflare.com/turnstile/v0/siteverify", url.Values{"secret": {config.Instance.TurnstileSecretKey}, "response": {gjson.Parse(util.BytesToString(ctx.Body())).Get("turnstileToken").String()}})
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		result := gjson.Parse(util.BytesToString(body))
		if !result.Get("success").Bool() {
			return response.New("failed", fiber.Map{"ec": result.Get("error-codes").String()}).Write(ctx, fiber.StatusBadRequest)
		}

		age := 24 * time.Hour
		exp := time.Now().Add(age).Unix()

		tokenStr, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"exp": exp,
		}).SignedString(util.StringToBytes(config.Instance.JWTSecret))
		if err != nil {
			return err
		}
		ctx.Cookie(&fiber.Cookie{
			Name:     "cfv_c",
			Value:    tokenStr,
			HTTPOnly: true,
			MaxAge:   int(age.Seconds()),
			Path:     "/",
		})
		return response.New("success").Write(ctx)
	}
}
