package metrics

import (
    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "time"
)

var (
    // Total requests counter
    httpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "route", "status"},
    )

    // Request duration histogram
    httpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "route"},
    )
)

// PrometheusMiddleware is a Gin middleware to collect HTTP metrics
func PrometheusMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        // Process request
        c.Next()
        
        // Record metrics after request is processed
        duration := time.Since(start).Seconds()
        status := c.Writer.Status()
        
        // Extract route pattern - in a real app, you might want to get the route pattern instead of the exact path
        route := c.Request.URL.Path
        method := c.Request.Method
        
        // Increment the requests counter
        httpRequestsTotal.WithLabelValues(method, route, string(rune(status))).Inc()
        
        // Observe the request duration
        httpRequestDuration.WithLabelValues(method, route).Observe(duration)
    }
}

// SetupPrometheusEndpoint adds the /metrics endpoint to the Gin router
func SetupPrometheusEndpoint(router *gin.Engine) {
    router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}