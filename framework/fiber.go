package framework

import (
	"time"

	"github.com/Sn0wo2/OpenCloudflareCDN/config"
	"github.com/Sn0wo2/OpenCloudflareCDN/debug"
	"github.com/Sn0wo2/OpenCloudflareCDN/router/errorhandler"
	"github.com/gofiber/fiber/v2"
)

func Fiber() *fiber.App {
	return fiber.New(fiber.Config{
		AppName:               "OpenCloudflareCDN",
		CaseSensitive:         true,
		DisableStartupMessage: false,
		ErrorHandler:          errorhandler.Error,
		IdleTimeout:           5 * time.Second,
		// dlv cant debug multiple process
		Prefork:           !debug.IsDebugging(),
		ReadTimeout:       10 * time.Second,
		ReduceMemoryUsage: true,
		StrictRouting:     true,
		WriteTimeout:      10 * time.Second,
		ServerHeader:      config.Instance.Server.Header,
	})
}

func Start(app *fiber.App) error {
	if config.Instance.Server.TLS.Cert != "" && config.Instance.Server.TLS.Key != "" {
		return app.ListenTLS(config.Instance.Server.Address, config.Instance.Server.TLS.Cert, config.Instance.Server.TLS.Key)
	}

	return app.Listen(config.Instance.Server.Address)
}
