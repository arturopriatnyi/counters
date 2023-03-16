package http

import (
	"encoding/json"
	"net/http"
	"time"

	"counters/pkg/counter"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type addCounterRequest struct {
	ID string `json:"id"`
}

func addCounter(l *zap.Logger, cm CounterManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		var r addCounterRequest
		if err := c.BindJSON(&r); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := cm.Add(r.ID)

		switch err {
		case nil:
			c.AbortWithStatus(http.StatusCreated)

			defer addCounterRequestDurationHistogram.With(nil).Observe(time.Since(start).Seconds())
			defer countersNumberGauge.With(nil).Inc()
		case counter.ErrExists:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			body, _ := json.Marshal(r)
			l.Error(
				"internal server error",
				zap.String("uri", c.Request.RequestURI),
				zap.String("body", string(body)),
				zap.Error(err),
			)

			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}

type getCounterResponse struct {
	ID    string `json:"id"`
	Value uint64 `json:"value"`
}

func getCounter(l *zap.Logger, cm CounterManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		c, err := cm.Get(id)

		switch err {
		case nil:
			ctx.AbortWithStatusJSON(http.StatusOK, getCounterResponse{
				ID:    c.ID,
				Value: c.Value,
			})
		case counter.ErrNotFound:
			ctx.AbortWithStatus(http.StatusNotFound)
		default:
			l.Error(
				"internal server error",
				zap.String("uri", ctx.Request.RequestURI),
				zap.String("id", id),
				zap.Error(err),
			)

			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}

func incCounter(l *zap.Logger, cm CounterManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		err := cm.Inc(id)

		switch err {
		case nil:
			ctx.AbortWithStatus(http.StatusOK)

			defer incCounterCounter.With(nil).Inc()
		case counter.ErrNotFound:
			ctx.AbortWithStatus(http.StatusNotFound)
		default:
			l.Error(
				"internal server error",
				zap.String("uri", ctx.Request.RequestURI),
				zap.String("id", id),
				zap.Error(err),
			)

			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}

func deleteCounter(l *zap.Logger, cm CounterManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		err := cm.Delete(id)

		switch err {
		case nil:
			ctx.AbortWithStatus(http.StatusNoContent)

			defer countersNumberGauge.With(nil).Dec()
		case counter.ErrNotFound:
			ctx.AbortWithStatus(http.StatusNotFound)
		default:
			l.Error(
				"internal server error",
				zap.String("uri", ctx.Request.RequestURI),
				zap.String("id", id),
				zap.Error(err),
			)

			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}
