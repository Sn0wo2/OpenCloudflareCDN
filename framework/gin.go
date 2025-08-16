package framework

import (
	"net/http"
	"time"

	"github.com/Sn0wo2/OpenCloudflareCDN/config"
	"github.com/Sn0wo2/OpenCloudflareCDN/debug"
	"github.com/Sn0wo2/OpenCloudflareCDN/log"
	"github.com/gin-gonic/gin"
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

// ZapLogger 使用 zap 记录生产环境的日志
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

	if config.Instance.Server.TLS.Cert != "" && config.Instance.Server.TLS.Key != "" {
		// http2 and http1
		g.Go(func() error {
			server := http.Server{
				Addr:    config.Instance.Server.Address,
				Handler: engine,
			}
			return server.ListenAndServeTLS(config.Instance.Server.TLS.Cert, config.Instance.Server.TLS.Key)
		})
		// http3
		g.Go(func() error {
			return http3.ListenAndServeQUIC(config.Instance.Server.Address, config.Instance.Server.TLS.Cert, config.Instance.Server.TLS.Key, engine)
		})
	} else {
		g.Go(func() error {
			return engine.Run(config.Instance.Server.Address)
		})
	}

	return g.Wait()
}
