package router

import (
	"github.com/Ai-feier/rbacapp/handler"
	ginmdl "github.com/Ai-feier/rbacapp/middleware/gin"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func NewRouter() *gin.Engine {
	router := gin.Default()

	// prometheus
	prom := handler.NewPrometheusHandler(promhttp.Handler())
	router.GET("/metrics", prom.Build())

	// prometheus gin route count middleware
	promGinCntMdl := ginmdl.NewPrometheusCountMiddleware(
		"rbacapi", "middleware",
		"gin_count", "the count of each request", map[string]string{
			"enviroment": "dev",
		})
	router.Use(promGinCntMdl.Build())

	// prometheus gin route summary middleware
	promGinSummaryMdl := ginmdl.NewPrometheusCountMiddleware(
		"rbacapi", "middleware",
		"gin_summary", "the time summary of each request", map[string]string{
			"enviroment": "dev",
		})
	router.Use(promGinSummaryMdl.Build())

	// prometheus gin route histogram middleware
	promGinHistogramMdl := ginmdl.NewPrometheusCountMiddleware(
		"rbacapi", "middleware",
		"gin_histogram", "the time histogram of each request", map[string]string{
			"enviroment": "dev",
		})
	router.Use(promGinHistogramMdl.Build())









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
