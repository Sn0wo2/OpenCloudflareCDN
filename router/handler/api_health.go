package handler

import (
	"github.com/Sn0wo2/OpenCloudflareCDN/internal/util"
	"github.com/Sn0wo2/OpenCloudflareCDN/log"
	"github.com/Sn0wo2/OpenCloudflareCDN/response"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func APIHealth() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		log.Instance.Info("H >> Health", zap.String("ctx", util.FiberContextString(ctx)))

		return response.New("ok").Write(ctx)
	}
}
