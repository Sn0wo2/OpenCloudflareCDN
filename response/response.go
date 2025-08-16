package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

func New(msg string, data ...any) *Response {
	switch len(data) {
	case 0:
		return &Response{
			Msg: msg,
		}
	case 1:
		return &Response{
			Msg:  msg,
			Data: data,
		}
	default:
		return &Response{
			Msg:  msg,
			Data: data,
		}
	}
}

func (r *Response) Write(ctx *gin.Context, status ...int) {
	s := http.StatusOK
	if len(status) > 0 {
		s = status[0]
	}

	ctx.JSON(s, r)
}
