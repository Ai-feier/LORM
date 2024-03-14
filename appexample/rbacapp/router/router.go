package router

import (
	"github.com/Ai-feier/rbacapp/handler"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter() *gin.Engine {
	router := gin.Default()

	user := router.Group("/api/user")
	userHandler := handler.NewUserHandler()
	{
		user.POST("/login", userHandler.LoginUser())
		user.POST("/registry", userHandler.Registry())
	}


	roleHandler := handler.NewRoleHandler()
	role := router.Group("/api/role")
	{
		role.GET("/ping", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, "success")
		})

		role.POST("/create", roleHandler.CreateRole())

		role.DELETE("/delete", roleHandler.DeleteRole())
	}


	rolebindingHandler := handler.NewRoleBindingHandler()
	rolebinding := router.Group("/api/rolebinding")
	{
		rolebinding.POST("/create", rolebindingHandler.CreateRoleBinding())
		rolebinding.DELETE("/delete", rolebindingHandler.DeleteRolebinding())
	}

	//clusterrole := router.Group("/api/clusterrole")
	//{
	//
	//}


	return router

}
