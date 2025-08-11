package main

import (
	"errors"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/Sn0wo2/OpenCloudflareCDN/config"
	"github.com/Sn0wo2/OpenCloudflareCDN/config/file"
	"github.com/Sn0wo2/OpenCloudflareCDN/framework"
	"github.com/Sn0wo2/OpenCloudflareCDN/log"
	"github.com/Sn0wo2/OpenCloudflareCDN/router"
	"github.com/Sn0wo2/OpenCloudflareCDN/version"
	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func init() {
	debug.SetGCPercent(50)

	_ = godotenv.Load()

	if err := config.Init(file.NewYAMLLoader(), file.NewJSONLoader()); err != nil {
		if errors.Is(err, config.ErrConfigNotFound) {
			envPath := config.Instance.ConfigPath
			if envPath == "" {
				envPath = "./data/config.yml"
			}

			config.Instance = config.DefaultConfig
			if err := file.NewYAMLLoader().Save(config.DefaultConfig, envPath); err != nil {
				panic(err)
			}

			config.Instance.IsDefault = true
			config.Instance.ConfigPath = envPath
		} else {
			panic(err)
		}
	}

	log.Init()
}

func main() {
	defer func() {
		_ = log.Instance.Sync()
	}()

	if !fiber.IsChild() {
		log.Instance.Info("Starting OpenCloudflareCDN...", zap.String("version", version.GetFormatVersion()))
	}

	if config.Instance.IsDefault {
		log.Instance.Warn("You have not set a configuration file, using the default!", zap.String("path", config.Instance.ConfigPath))
	}

	app := framework.Fiber()

	router.Init(app)

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := framework.Start(app); err != nil {
			log.Instance.Fatal("Server failed to start",
				zap.Error(err),
			)
		}
	}()

	<-shutdownChan

	if err := app.Shutdown(); err != nil {
		log.Instance.Error("Server failed to shutdown",
			zap.Error(err),
		)
	}
}
