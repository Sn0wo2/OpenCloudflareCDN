package router

import (
	"github.com/Sn0wo2/OpenCloudflareCDN/router/handler"
	"github.com/Sn0wo2/OpenCloudflareCDN/router/notfound"
	"github.com/Sn0wo2/OpenCloudflareCDN/router/static"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Init(router fiber.Router) {
	router.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}), cors.New())

	debug := router.Group("/v0")
	debug.Get("/error", handler.APIError())

	api := router.Group("/v1")
	api.Get("/health", handler.APIHealth())

	api.Post("/verify", handler.APIVerify())

	router.All("*", handler.Proxy())

	static.Init(router)
	notfound.Init(router)
}
