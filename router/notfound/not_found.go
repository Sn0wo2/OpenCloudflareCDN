package notfound

import (
	"strings"

	"github.com/Sn0wo2/OpenCloudflareCDN/internal/util"
	"github.com/Sn0wo2/OpenCloudflareCDN/log"
	"github.com/Sn0wo2/OpenCloudflareCDN/response"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Init(router fiber.Router, msg ...string) {
	m := strings.Join(msg, " ")
	if m = strings.ToLower(m); m == "" {
		m = "page not found"
	}

	router.Use("*", func(ctx *fiber.Ctx) error {
		log.Instance.Warn("NF >> "+util.TitleCase(m),
			zap.String("ctx", util.FiberContextString(ctx)))

		return response.New(m).Write(ctx, fiber.StatusNotFound)
	})
}
