package handler

import (
	"github.com/Ai-feier/rbacapp/pkg/resp"
	"github.com/Ai-feier/rbacapp/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RoleBindingHandler struct {
	roleBindingSvc *service.RoleBindingSvc
}

func NewRoleBindingHandler() *RoleBindingHandler {
	return &RoleBindingHandler{
		roleBindingSvc: service.NewRoleBindingSvc(),
	}
}

func (r RoleBindingHandler) CreateRoleBinding() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var rb RoleBindingRequest
		if err := ctx.Bind(&rb); err != nil {
			ctx.JSON(http.StatusBadRequest, resp.RespError(ctx, err, "参数绑定失败"))
			return
		}

		err := r.roleBindingSvc.CreateRoleBinding(ctx, rb.RoleBinding, rb.Role)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, resp.RespError(ctx, err, "rolebinding 创建失败"))
			return
		}
		ctx.JSON(http.StatusOK, resp.RespSuccess(ctx, "rolebinding 创建成功"))
	}
}

func (r RoleBindingHandler) DeleteRolebinding() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var rb RoleBindingRequest
		if err := ctx.Bind(&rb); err != nil {
			ctx.JSON(http.StatusBadRequest, resp.RespError(ctx, err, "参数绑定失败"))
			return
		}

		err := r.roleBindingSvc.DeleteRoleBinding(ctx, rb.RoleBinding)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, resp.RespError(ctx, err, "rolebinding 删除失败"))
			return
		}
		ctx.JSON(http.StatusOK, resp.RespSuccess(ctx, "rolebinding 删除成功"))
	}
}
