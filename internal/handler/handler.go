//go:generate mockgen -source=handler.go -destination=mock.go -package=handler
package handler

import (
	"context"
	"net/http"

	"counters/pkg/counter"
	"counters/pkg/oauth2"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type IAManager interface {
	OAuth2URL(provider oauth2.Provider) (string, error)
	SignInWithOAuth2(ctx context.Context, provider oauth2.Provider, state, code string) (oauth2.Token, error)
}

type CounterManager interface {
	Add(id string) error
	Get(id string) (*counter.Counter, error)
	Inc(id string) error
	Delete(id string) error
}

func New(l *zap.Logger, iam IAManager, cm CounterManager) http.Handler {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), withInternalServerErrorCounter())
	r.HandleMethodNotAllowed = true
	r.NoRoute(noRoute())
	r.NoMethod(noMethod())

	r.GET("/health", health())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	oauth := r.Group("/oauth")
	google := oauth.Group("/google")
	google.GET("/sign-in", googleSignIn(l, iam))
	google.GET("/callback", googleCallback(l, iam))
	github := oauth.Group("/github")
	github.GET("/sign-in", gitHubSignIn(l, iam))
	github.GET("/callback", githubCallback(l, iam))

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
