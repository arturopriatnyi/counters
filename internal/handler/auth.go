package handler

import (
	"net/http"

	"counters/pkg/oauth2"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func googleSignIn(l *zap.Logger, iam IAManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		url, err := iam.OAuth2URL(oauth2.Google)
		if err != nil {
			l.Error(
				"internal server error",
				zap.String("uri", c.Request.RequestURI),
				zap.Error(err),
			)

			c.Status(http.StatusInternalServerError)
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

type googleCallbackResponse struct {
	AccessToken string          `json:"access_token"`
	Provider    oauth2.Provider `json:"provider"`
}

func googleCallback(l *zap.Logger, iamManager IAManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		state, code := c.Query("state"), c.Query("code")

		token, err := iamManager.SignInWithOAuth2(c, oauth2.Google, state, code)
		if err != nil {
			l.Error(
				"internal server error",
				zap.String("uri", c.Request.RequestURI),
				zap.Error(err),
			)

			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, googleCallbackResponse{
			AccessToken: token.Access,
			Provider:    token.Provider,
		})
	}
}

func gitHubSignIn(l *zap.Logger, iam IAManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		url, err := iam.OAuth2URL(oauth2.GitHub)
		if err != nil {
			l.Error(
				"internal server error",
				zap.String("uri", c.Request.RequestURI),
				zap.Error(err),
			)

			c.Status(http.StatusInternalServerError)
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

type githubCallbackResponse struct {
	AccessToken string          `json:"access_token"`
	Provider    oauth2.Provider `json:"provider"`
}

func githubCallback(l *zap.Logger, iamManager IAManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		state, code := c.Query("state"), c.Query("code")

		token, err := iamManager.SignInWithOAuth2(c, oauth2.GitHub, state, code)
		if err != nil {
			l.Error(
				"internal server error",
				zap.String("uri", c.Request.RequestURI),
				zap.Error(err),
			)

			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, githubCallbackResponse{
			AccessToken: token.Access,
			Provider:    token.Provider,
		})
	}
}
