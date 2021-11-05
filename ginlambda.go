package ginlambda

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gin-gonic/gin"
)

// HandlerFunc is the signature for the Lambda handler function.
type HandlerFunc func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

// Start is analogous to lambda.Start() but takes a *gin.Engine argument instead of a handler
// function. The engine should have any desired routes initialized but should not be run.
func Start(r *gin.Engine) {
	handler := NewHandler(r)
	lambda.Start(handler)
}

// NewHandler creates a new Lambda handler function from a *gin.Engine instance. This handler may
// be passed as the handler argument to lambda.Start().
func NewHandler(r *gin.Engine) HandlerFunc {
	return func(ctx context.Context,
		request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

		httpRequest, useMultiValueHeader, err := constructHTTPRequestFromAPIRequest(ctx, request)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		collector := newAPIResponseCollector(useMultiValueHeader)

		r.ServeHTTP(collector, httpRequest)

		return collector.ToAPIGatewayProxyResponse(), nil
	}
}

func constructHTTPRequestFromAPIRequest(
	ctx context.Context, request events.APIGatewayProxyRequest) (*http.Request, bool, error) {

	useMultiValueHeader := (request.Headers == nil)

	// initialize request query
	var queryStr string
	if useMultiValueHeader {
		queryStr = url.Values(request.MultiValueQueryStringParameters).Encode()
	} else {
		queryValues := make(url.Values)
		for k, v := range request.QueryStringParameters {
			queryValues.Set(k, v)
		}
		queryStr = queryValues.Encode()
	}

	// initialize request with context
	reader := ioutil.NopCloser(strings.NewReader(request.Body))
	fullPath := fmt.Sprintf("%s?%s", request.Path, queryStr)
	httpRequest, err := http.NewRequestWithContext(ctx, request.HTTPMethod, fullPath, reader)
	if err != nil {
		return nil, false, err
	}

	// initialize request header
	if useMultiValueHeader {
		for k, vs := range request.MultiValueHeaders {
			for _, v := range vs {
				httpRequest.Header.Add(k, v)
			}
		}
	} else {
		for k, v := range request.Headers {
			httpRequest.Header.Set(k, v)
		}
	}

	return httpRequest, useMultiValueHeader, nil
}
