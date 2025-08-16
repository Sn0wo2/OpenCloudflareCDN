package notfound

import (
	"net/http"
	"strings"

	"github.com/Sn0wo2/OpenCloudflareCDN/internal/util"
	"github.com/Sn0wo2/OpenCloudflareCDN/log"
	"github.com/Sn0wo2/OpenCloudflareCDN/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Handler(msg ...string) gin.HandlerFunc {
	m := strings.Join(msg, " ")
	if m = strings.ToLower(m); m == "" {
		m = "page not found"
	}

	return func(ctx *gin.Context) {
		log.Instance.Warn("NF >> "+util.TitleCase(m),
			zap.String("ctx", util.GinContextString(ctx)))

		response.New(m).Write(ctx, http.StatusNotFound)
	}
}
