package errorhandler

import (
	"net/http"

	"github.com/Sn0wo2/OpenCloudflareCDN/internal/util"
	"github.com/Sn0wo2/OpenCloudflareCDN/log"
	"github.com/Sn0wo2/OpenCloudflareCDN/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func Error() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) > 0 {
			for _, err := range ctx.Errors {
				traceID := uuid.NewString()

				log.Instance.Error("EH >> Error handler caught error",
					zap.String("traceID", traceID),
					zap.Error(err.Err),
					zap.String("ctx", util.GinContextString(ctx)),
				)

				response.New("oops, something went wrong", gin.H{"traceID": traceID}).Write(ctx, http.StatusInternalServerError)
			}
		}
	}
}
