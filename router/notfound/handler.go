package notfound

import (
	"net/http"

	"github.com/Sn0wo2/OpenCloudflareCDN/internal/util"
	"github.com/Sn0wo2/OpenCloudflareCDN/log"
	"github.com/Sn0wo2/OpenCloudflareCDN/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Handler(msg ...string) gin.HandlerFunc {
	m := util.GetStatusText(http.StatusNotFound, msg...)

	return func(ctx *gin.Context) {
		log.Instance.Warn("NF >> "+util.TitleCase(m),
			zap.String("ctx", util.GinContextString(ctx)))

		response.New(m).Write(ctx, http.StatusMethodNotAllowed)
	}
}
