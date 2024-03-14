package resp

import (
	"github.com/Ai-feier/rbacapp/pkg/msg"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code  int    `json:"code,omitempty"`
	Data  any    `json:"data,omitempty"`
	Msg   string `json:"msg,omitempty"`
	Error string `json:"error,omitempty"`
}

func RespSuccess(ctx *gin.Context, data any, code ...int) *Response {
	status := msg.SUCCESS
	if code != nil {
		status = code[0]
	}

	if data == "" {
		data = "操作成功"
	}

	r := &Response{
		Code:  status,
		Data:  data,
		Msg:   msg.GetMsg(status),
	}
	return r
}

func RespError(ctx *gin.Context, err error, data string, code...int) *Response {
	status := msg.ERROR
	if code != nil {
		status = code[0]
	}

	r := &Response{
		Code:  status,
		Data:  data,
		Msg:   msg.GetMsg(status),
		Error: err.Error(),
	}
	return r
}
