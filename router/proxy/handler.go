package proxy

import (
	"fmt"

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
			log.Instance.Info("No token, not proxying",
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
			log.Instance.Error("Invalid token, not proxying",
				zap.String("token", tokenStr),
				zap.String("ctx", util.GinContextString(ctx)),
			)

			if err != nil {
				_ = ctx.Error(err)
			}
			ctx.Next()
			return
		}

		Request(ctx)
	}
}
