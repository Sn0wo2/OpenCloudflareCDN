package proxy

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/Sn0wo2/OpenCloudflareCDN/config"
	"github.com/Sn0wo2/OpenCloudflareCDN/internal/util"
	"github.com/Sn0wo2/OpenCloudflareCDN/log"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func Proxy() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr, err := ctx.Cookie("__ocfc_v")
		if err != nil {
			log.Instance.Info("P >> No token, not proxying",
				zap.String("ctx", util.GinContextString(ctx)),
				zap.Error(err),
			)
			ctx.Next()

			return
		}

		t, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return util.StringToBytes(config.Instance.JWTSecret), nil
		})
		if err != nil || !t.Valid {
			log.Instance.Error("P >> Invalid token, not proxying",
				zap.String("token", tokenStr),
				zap.String("ctx", util.GinContextString(ctx)),
			)

			if err != nil {
				_ = ctx.Error(err)
			}

			ctx.Next()

			return
		}

		claims, ok := t.Claims.(jwt.MapClaims)
		if !ok {
			log.Instance.Error("P >> Invalid token claims, not proxying", zap.String("ctx", util.GinContextString(ctx)))
			ctx.Next()

			return
		}

		clientIP, ok1 := claims["ip"].(string)

		userAgent, ok2 := claims["ua"].(string)
		if !ok1 || !ok2 || clientIP != ctx.ClientIP() || userAgent != ctx.Request.UserAgent() {
			log.Instance.Warn(
				"P >> Token information mismatch",
				zap.String("tokenIP", clientIP),
				zap.String("requestIP", ctx.ClientIP()),
				zap.String("tokenUA", userAgent),
				zap.String("requestUA", ctx.Request.UserAgent()),
				zap.String("ctx", util.GinContextString(ctx)),
			)
			ctx.Next()

			return
		}

		remote, err := url.Parse(config.Instance.OriginalServer)
		if err != nil {
			_ = ctx.Error(err)

			return
		}

		log.Instance.Info("P >> Proxying request", zap.String("target", remote.String()), zap.String("protocol", "HTTP/2"), zap.String("ctx", util.GinContextString(ctx)))

		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec
		}
		proxy.Director = func(req *http.Request) {
			req.Header = ctx.Request.Header
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = ctx.Request.URL.Path
		}
		proxy.ServeHTTP(ctx.Writer, ctx.Request)
	}
}
