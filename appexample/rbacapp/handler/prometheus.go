package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type PrometheusHandler struct{
	handler http.Handler
}

func NewPrometheusHandler(handler http.Handler) *PrometheusHandler {
	return &PrometheusHandler{
		handler: handler,
	}
}

func (p PrometheusHandler) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		p.handler.ServeHTTP(ctx.Writer, ctx.Request)
	}
}
