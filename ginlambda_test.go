package ginlambda_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/dgravesa/ginlambda"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func unmarshalJSONString(t *testing.T, str string, v interface{}) {
	err := json.Unmarshal([]byte(str), v)
	if err != nil {
		t.Error("error on unmarshal json:", err)
	}
}

func Test_HandlerFunc_WithRouteParams_ReturnsExpectedResult(t *testing.T) {
	// arrange
	type responseBody struct {
		UserID   string `json:"userId"`
		Greeting string `json:"greeting"`
	}
	// create gin engine with handler
	r := gin.New()
	r.GET("/users/:userID/greetings/:greetingID", func(c *gin.Context) {
		userID := c.Param("userID")
		greetingID := c.Param("greetingID")

		greetingFmt, found := map[string]string{
			"morning":   "Good morning, %s!",
			"afternoon": "Great afternoon, %s!",
			"evening":   "Fair evening, %s!",
		}[greetingID]

		if !found {
			c.Status(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, responseBody{
			UserID:   userID,
			Greeting: fmt.Sprintf(greetingFmt, userID),
		})
	})
	handler := ginlambda.NewHandler(r)
	// construct test request
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Path:       "/users/Dave/greetings/afternoon",
	}
	expectedResponse := responseBody{
		UserID:   "Dave",
		Greeting: "Great afternoon, Dave!",
	}

	// act
	response, err := handler(context.Background(), request)

	// assert
	if err != nil {
		t.Fatal("expected nil, received error:", err)
	}
	var actualResponse responseBody
	unmarshalJSONString(t, response.Body, &actualResponse)
	if expectedResponse != actualResponse {
		t.Errorf("expected: %v, received: %v\n", expectedResponse, actualResponse)
	}
	if http.StatusOK != response.StatusCode {
		t.Errorf("expected status code: %d, received status code: %d\n", http.StatusOK, response.StatusCode)
	}
}

func Test_HandlerFunc_WithSingleValueQueryParams_ReturnsExpectedResult(t *testing.T) {
	// arrange
	type responseBody struct {
		Greeting string `json:"greeting"`
	}
	// create gin engine with handler
	r := gin.New()
	r.GET("/greeting", func(c *gin.Context) {
		name := c.Query("name")
		timeOfDay := c.Query("timeOfDay")

		greetingFmt, found := map[string]string{
			"morning":   "Good morning, %s!",
			"afternoon": "Great afternoon, %s!",
			"evening":   "Fair evening, %s!",
		}[timeOfDay]

		if !found {
			c.Status(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, responseBody{
			Greeting: fmt.Sprintf(greetingFmt, name),
		})
	})
	handler := ginlambda.NewHandler(r)
	// construct test request
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Path:       "/greeting",
		QueryStringParameters: map[string]string{
			"name":      "Spud",
			"timeOfDay": "morning",
		},
		Headers: map[string]string{},
	}
	expectedResponse := responseBody{
		Greeting: "Good morning, Spud!",
	}

	// act
	response, err := handler(context.Background(), request)

	// assert
	if err != nil {
		t.Fatal("expected nil, received error:", err)
	}
	var actualResponse responseBody
	unmarshalJSONString(t, response.Body, &actualResponse)
	if expectedResponse != actualResponse {
		t.Errorf("expected: %v, received: %v\n", expectedResponse, actualResponse)
	}
	if http.StatusOK != response.StatusCode {
		t.Errorf("expected status code: %d, received status code: %d\n", http.StatusOK, response.StatusCode)
	}
}

func Test_HandlerFunc_WithMultiValueQueryParams_ReturnsExpectedResult(t *testing.T) {
	// arrange
	type responseBody struct {
		Greeting string `json:"greeting"`
	}
	// create gin engine with handler
	r := gin.New()
	r.GET("/greeting", func(c *gin.Context) {
		names := c.QueryArray("name")
		timeOfDay := c.Query("timeOfDay")

		greetingFmt, found := map[string]string{
			"morning":   "Good morning, %s!",
			"afternoon": "Great afternoon, %s!",
			"evening":   "Fair evening, %s!",
		}[timeOfDay]

		if !found {
			c.Status(http.StatusNotFound)
			return
		}

		namesStr := strings.Join(names, " and ")
		c.JSON(http.StatusOK, responseBody{
			Greeting: fmt.Sprintf(greetingFmt, namesStr),
		})
	})
	handler := ginlambda.NewHandler(r)
	// construct test request
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Path:       "/greeting",
		MultiValueQueryStringParameters: map[string][]string{
			"name":      {"Trent", "Sydney"},
			"timeOfDay": {"evening"},
		},
		MultiValueHeaders: map[string][]string{},
	}
	expectedResponse := responseBody{
		Greeting: "Fair evening, Trent and Sydney!",
	}

	// act
	response, err := handler(context.Background(), request)

	// assert
	if err != nil {
		t.Fatal("expected nil, received error:", err)
	}
	var actualResponse responseBody
	unmarshalJSONString(t, response.Body, &actualResponse)
	if expectedResponse != actualResponse {
		t.Errorf("expected: %v, received: %v\n", expectedResponse, actualResponse)
	}
	if http.StatusOK != response.StatusCode {
		t.Errorf("expected status code: %d, received status code: %d\n", http.StatusOK, response.StatusCode)
	}
}

func Test_HandlerFunc_WithRequestBody_ReturnsExpectedResult(t *testing.T) {
	// arrange

	// act

	// assert
}

func Test_HandlerFunc_WithSingleValueHeaders_ReturnsExpectedResult(t *testing.T) {
	// arrange

	// act

	// assert
}

func Test_HandlerFunc_WithMultiValueHeaders_ReturnsExpectedResult(t *testing.T) {
	// arrange

	// act

	// assert
}
