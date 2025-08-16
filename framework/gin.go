package framework

import (
	"net/http"
	"time"

	"github.com/Sn0wo2/OpenCloudflareCDN/config"
	"github.com/Sn0wo2/OpenCloudflareCDN/debug"
	"github.com/Sn0wo2/OpenCloudflareCDN/log"
	"github.com/gin-gonic/gin"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func Gin() *gin.Engine {
	if !debug.IsDebugging() {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()

	engine.Use(gin.Recovery())

	engine.Use(ZapLogger())

	engine.HandleMethodNotAllowed = true

	return engine
}

func ZapLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery

		ctx.Next()

		if raw != "" {
			path = path + "?" + raw
		}

		log.Instance.Info("GIN",
			zap.Int("status", ctx.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
			zap.String("clientIP", ctx.ClientIP()),
			zap.String("method", ctx.Request.Method),
			zap.String("path", path),
			zap.String("error", ctx.Errors.ByType(gin.ErrorTypePrivate).String()),
		)
	}
}

func Start(engine *gin.Engine) error {
	g := new(errgroup.Group)

	server := &http.Server{
		Addr:              config.Instance.Server.Address,
		Handler:           engine,
		ReadHeaderTimeout: 30 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	if config.Instance.Server.TLS.Cert != "" && config.Instance.Server.TLS.Key != "" {
		// http2 and http1
		g.Go(func() error {
			return server.ListenAndServeTLS(config.Instance.Server.TLS.Cert, config.Instance.Server.TLS.Key)
		})
		// http3
		g.Go(func() error {
			h3Server := &http3.Server{
				Addr:    config.Instance.Server.Address,
				Handler: engine,
				QUICConfig: &quic.Config{
					MaxIdleTimeout:       30 * time.Second,
					HandshakeIdleTimeout: 10 * time.Second,
				},
			}

			return h3Server.ListenAndServeTLS(config.Instance.Server.TLS.Cert, config.Instance.Server.TLS.Key)
		})
	} else {
		g.Go(func() error {
			return server.ListenAndServe()
		})
	}

	return g.Wait()
}
