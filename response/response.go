package response

import (
	"github.com/gofiber/fiber/v2"
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
			Data: data[0],
		}
	default:
		return &Response{
			Msg:  msg,
			Data: data,
		}
	}
}

func (r *Response) Write(ctx *fiber.Ctx, status ...int) error {
	if len(status) > 0 {
		ctx.Status(status[0])
	}

	return ctx.JSON(r)
}
