package metrics

import (
	"strconv"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    httpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "route", "status"},
    )

    httpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "route"},
    )
)

func PrometheusMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start).Seconds()
        status := c.Writer.Status()
        
        route := c.Request.URL.Path
        method := c.Request.Method
        
        httpRequestsTotal.WithLabelValues(method, route, strconv.Itoa(status)).Inc()
        
        httpRequestDuration.WithLabelValues(method, route).Observe(duration)
    }
}

func SetupPrometheusEndpoint(router *gin.Engine) {
    router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}