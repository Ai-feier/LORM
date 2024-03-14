package handler

import (
	"context"
	"github.com/Ai-feier/rbacapp/pkg/resp"
	"github.com/Ai-feier/rbacapp/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct{
}

func NewUserHandler() UserHandler {
	return UserHandler{}
}


// LoginUser 用户身份认证
func (u UserHandler) LoginUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var userReq UserRequest
		if err := ctx.Bind(&userReq); err != nil {
			ctx.JSON(http.StatusBadRequest, resp.RespError(ctx, err, "参数绑定失败"))
			return
		}

		userSvc := service.NewUserSvc()
		user, err := userSvc.Login(ctx, userReq.User.UserName, userReq.User.Password)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, resp.RespError(ctx, err, "登录失败"))
			return
		}
		ctx.JSON(http.StatusOK, resp.RespSuccess(ctx, user))
	}
}

func (u UserHandler) Registry() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var userReq UserRequest
		if err := ctx.Bind(&userReq); err != nil {
			ctx.JSON(http.StatusBadRequest, resp.RespError(ctx, err, "参数绑定失败"))
			return
		}

		userSvc := service.NewUserSvc()
		user, err := userSvc.Registry(context.Background(), userReq.User.UserName, userReq.User.Password)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, resp.RespError(ctx, err, "注册失败"))
		}
		ctx.JSON(http.StatusOK, resp.RespSuccess(ctx, user))
	}
}
