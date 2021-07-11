package ginlambda_test

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/dgravesa/ginlambda"
	"github.com/gin-gonic/gin"
)

func ExampleNewHandler() {
	// initialize gin route
	r := gin.Default()
	r.GET("/greeting/:userName", func(c *gin.Context) {
		userName := c.Param("userName")
		c.String(http.StatusOK, "Hello, %s!", userName)
	})

	// test request
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Path:       "/greeting/Bruce",
	}

	// construct lambda handler from gin engine
	handler := ginlambda.NewHandler(r)

	// execute request
	response, _ := handler(context.Background(), request)

	fmt.Println(response.StatusCode, response.Body)
	// Output: 200 Hello, Bruce!
}
