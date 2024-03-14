package handler

import (
	"context"
	"github.com/Ai-feier/rbacapp/pkg/resp"
	"github.com/Ai-feier/rbacapp/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RoleHandler struct {
	roleSvc *service.RoleSvc
}

func NewRoleHandler() *RoleHandler {
	return &RoleHandler{
		roleSvc: service.NewRoleSvc(),
	}
}

func (r RoleHandler) CreateRole() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var roleReq RoleRequest
		if err := ctx.Bind(&roleReq); err != nil {
			ctx.JSON(http.StatusBadRequest, resp.RespError(ctx, err, "参数绑定失败"))
			return
		}

		err := r.roleSvc.CreateRole(context.Background(), roleReq.Role, roleReq.RoleSubRefs...)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, resp.RespError(ctx, err, "role 创建失败"))
			return
		}
		ctx.JSON(http.StatusOK, resp.RespSuccess(ctx, "role 创建成功"))
	}
}

func (r RoleHandler) DeleteRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RoleRequest
		if err := c.Bind(&req); err != nil {
			c.JSON(http.StatusBadRequest, resp.RespError(c, err, "参数绑定失败"))
			return
		}

		err := r.roleSvc.DeleteRole(c, req.Role)
		if err != nil {
			c.JSON(http.StatusBadRequest, resp.RespError(c, err, "role 删除失败"))
			return
		}
		c.JSON(http.StatusOK, resp.RespSuccess(c, "role 删除成功"))
	}
}
