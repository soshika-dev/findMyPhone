package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

// Setup attaches common middleware to gin Engine.
func Setup(r *gin.Engine) {
	r.Use(gin.Recovery())
	r.Use(requestid.New())
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
		ExposeHeaders:   []string{"X-Request-ID"},
		MaxAge:          12 * time.Hour,
	}))
}
