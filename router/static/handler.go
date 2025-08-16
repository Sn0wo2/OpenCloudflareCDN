package static

import (
	"net/http"

	"github.com/Sn0wo2/OpenCloudflareCDN/config"
	"github.com/gin-gonic/gin"
)

func Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method == http.MethodGet {
			staticRoot := config.Instance.StaticPath
			if staticRoot == "" {
				ctx.Next()

				return
			}

			filePath := GetStaticFile(config.Instance.StaticIndex, staticRoot, ctx.Request.URL.Path)
			if filePath == "" {
				ctx.Next()

				return
			}

			ctx.File(filePath)
			ctx.Abort()
		}
	}
}
