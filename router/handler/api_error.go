package handler

import (
	"errors"

	"github.com/Sn0wo2/OpenCloudflareCDN/internal/util"
	"github.com/Sn0wo2/OpenCloudflareCDN/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func APIError() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Instance.Info("E >> Error test", zap.String("ctx", util.GinContextString(ctx)))

		_ = ctx.Error(errors.New("test error"))
	}
}
