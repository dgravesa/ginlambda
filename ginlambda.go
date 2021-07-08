package ginlambda

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gin-gonic/gin"
)

type lambdaHandlerFunc func(events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error)

// Start is analogous to lambda.Start() but takes a *gin.Engine argument instead of a handler
// function. The engine should have any desired routes initialized but should not be run.
func Start(r *gin.Engine) {
	lambdaHandler := makeLambdaHandlerFromEngine(r)
	lambda.Start(lambdaHandler)
}

func makeLambdaHandlerFromEngine(r *gin.Engine) lambdaHandlerFunc {
	return func(request events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
		httpRequest, useMultiValueHeader, err := constructHTTPRequestFromALBRequest(request)
		if err != nil {
			return events.ALBTargetGroupResponse{}, err
		}

		collector := newALBResponseCollector(useMultiValueHeader)

		r.ServeHTTP(collector, httpRequest)

		return collector.ToALBTargetGroupResponse(), nil
	}
}

func constructHTTPRequestFromALBRequest(
	albRequest events.ALBTargetGroupRequest) (*http.Request, bool, error) {

	useMultiValueHeader := (albRequest.Headers == nil)

	// initialize request
	httpRequest := &http.Request{
		Method: albRequest.HTTPMethod,
		URL: &url.URL{
			Path: albRequest.Path,
		},
		Header: make(http.Header),
		Body:   ioutil.NopCloser(strings.NewReader(albRequest.Body)),
	}
	// manually add the raw path just in case
	httpRequest.URL.RawPath = httpRequest.URL.EscapedPath()

	// initialize request header
	if useMultiValueHeader {
		for k, vs := range albRequest.MultiValueHeaders {
			for _, v := range vs {
				httpRequest.Header.Add(k, v)
			}
		}
	} else {
		for k, v := range albRequest.Headers {
			httpRequest.Header.Set(k, v)
		}
	}

	// initialize request query
	if useMultiValueHeader {
		for k, vs := range albRequest.MultiValueQueryStringParameters {
			for _, v := range vs {
				httpRequest.URL.Query().Add(k, v)
			}
		}
	} else {
		for k, v := range albRequest.QueryStringParameters {
			httpRequest.URL.Query().Set(k, v)
		}
	}

	return httpRequest, useMultiValueHeader, nil
}
