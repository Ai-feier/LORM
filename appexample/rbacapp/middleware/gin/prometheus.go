package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

// PrometheusCountMiddleware 请求路由计数
type PrometheusCountMiddleware struct{
	Namespace string
	Subsystem string
	Name string
	Help string
	ConstLabels map[string]string
}

func NewPrometheusCountMiddleware(namespace, subsystem, name, help string, constLabels map[string]string) *PrometheusCountMiddleware {
	return &PrometheusCountMiddleware{
		Namespace:   namespace,
		Subsystem:   subsystem,
		Name:        name,
		Help:        help,
		ConstLabels: constLabels,
	}
}

func (p PrometheusCountMiddleware) Build() gin.HandlerFunc {
	count := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   p.Namespace,
		Subsystem:   p.Subsystem,
		Name:        p.Name,
		Help:        p.Help,
		ConstLabels: p.ConstLabels,
	}, []string{"pattern", "method", "status"})

	prometheus.MustRegister(count)

	return func(ctx *gin.Context) {
		defer func() {
			count.WithLabelValues(ctx.Request.URL.Path, ctx.Request.Method, strconv.Itoa(ctx.Writer.Status())).
				Add(1)
		}()
		ctx.Next()
	}
}

// PrometheusSummaryMiddleware 请求超时计数 99线 999线
type PrometheusSummaryMiddleware struct{
	Namespace string
	Subsystem string
	Name string
	Help string
	ConstLabels map[string]string
}

func NewPrometheusSummaryMiddleware(namespace, subsystem, name, help string, constLabels map[string]string) *PrometheusSummaryMiddleware {
	return &PrometheusSummaryMiddleware{
		Namespace:   namespace,
		Subsystem:   subsystem,
		Name:        name,
		Help:        help,
		ConstLabels: constLabels,
	}
}

func (p PrometheusSummaryMiddleware) Build() gin.HandlerFunc {
	summary := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:   p.Namespace,
		Subsystem:   p.Subsystem,
		Name:        p.Name,
		Help:        p.Help,
		ConstLabels: p.ConstLabels,
		Objectives: map[float64]float64{
			0.5: 0.01,
			0.75: 0.01,
			0.90: 0.01,
			0.99: 0.001,
			0.999: 0.0001,
		},
	}, []string{"pattern", "method", "status"})

	prometheus.MustRegister(summary)

	return func(ctx *gin.Context) {
		start := time.Now()
		defer func() {
			duration := time.Now().Sub(start).Milliseconds()
			summary.WithLabelValues(ctx.Request.URL.Path, ctx.Request.Method,
				strconv.Itoa(ctx.Writer.Status())).Observe(float64(duration))
		}()
		ctx.Next()
	}
}


// PrometheusHistogramMiddleware 请求超时计数 99线 999线
type PrometheusHistogramMiddleware struct{
	Namespace string
	Subsystem string
	Name string
	Help string
	ConstLabels map[string]string
}

func NewPrometheusHistogramMiddleware(namespace, subsystem, name, help string, constLabels map[string]string) *PrometheusHistogramMiddleware {
	return &PrometheusHistogramMiddleware{
		Namespace:   namespace,
		Subsystem:   subsystem,
		Name:        name,
		Help:        help,
		ConstLabels: constLabels,
	}
}

func (p PrometheusHistogramMiddleware) Build() gin.HandlerFunc {
	histogram  := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   p.Namespace,
		Subsystem:   p.Subsystem,
		Name:        p.Name,
		Help:        p.Help,
		ConstLabels: p.ConstLabels,
		Buckets: []float64{1, 3, 5, 7, 10, 17, 23, 30, 50, 100, 1000},
	}, []string{"pattern", "method", "status"})

	prometheus.MustRegister(histogram)

	return func(ctx *gin.Context) {
		start := time.Now()
		defer func() {
			duration := time.Now().Sub(start).Milliseconds()
			histogram.WithLabelValues(ctx.Request.URL.Path, ctx.Request.Method,
				strconv.Itoa(ctx.Writer.Status())).Observe(float64(duration))
		}()
		ctx.Next()
	}
}
