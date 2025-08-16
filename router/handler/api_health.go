package handler

import (
	"github.com/Sn0wo2/OpenCloudflareCDN/internal/util"
	"github.com/Sn0wo2/OpenCloudflareCDN/log"
	"github.com/Sn0wo2/OpenCloudflareCDN/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func APIHealth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Instance.Info("H >> Health", zap.String("ctx", util.GinContextString(ctx)))

		response.New("ok").Write(ctx)
	}
}
