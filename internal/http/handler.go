package http

//go:generate mockgen -source=handler.go -destination=mock.go -package=http

import (
	"net/http"

	"counters/pkg/counter"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type CounterManager interface {
	Add(id string) error
	Get(id string) (counter.Counter, error)
	Inc(id string) error
	Delete(id string) error
}

func NewHandler(l *zap.Logger, cm CounterManager) http.Handler {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), withInternalServerErrorCounter())
	r.HandleMethodNotAllowed = true

	r.NoRoute(noRoute())
	r.NoMethod(noMethod())
	r.GET("/health", health())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	counters := r.Group("/counters")
	counters.POST("", addCounter(l, cm))
	counters.GET("/:id", getCounter(l, cm))
	counters.GET("/:id/inc", incCounter(l, cm))
	counters.DELETE("/:id", deleteCounter(l, cm))

	return r
}

func noRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func noMethod() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatus(http.StatusMethodNotAllowed)
	}
}

func health() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatus(http.StatusOK)
	}
}
