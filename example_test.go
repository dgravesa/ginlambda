package ginlambda_test

import (
	"net/http"

	"github.com/dgravesa/ginlambda"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func init() {
	r = gin.Default()
	r.GET("/greeting", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})
}

func Example() {
	ginlambda.Start(r)
}
