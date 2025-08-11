package handler

import (
	"errors"

	"github.com/Sn0wo2/OpenCloudflareCDN/internal/util"
	"github.com/Sn0wo2/OpenCloudflareCDN/log"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func APIError() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		log.Instance.Info("E >> Error test", zap.String("ctx", util.FiberContextString(ctx)))

		return errors.New("test error")
	}
}
