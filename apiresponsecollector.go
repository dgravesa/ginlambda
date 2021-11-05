package ginlambda

import (
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type apiResponseCollector struct {
	header              map[string][]string
	stringbuilder       strings.Builder
	statusCode          int
	useMultiValueHeader bool
}

func newAPIResponseCollector(useMultiValueHeader bool) *apiResponseCollector {
	return &apiResponseCollector{
		header:              make(map[string][]string),
		useMultiValueHeader: useMultiValueHeader,
	}
}

func (rc *apiResponseCollector) Header() http.Header {
	return rc.header
}

func (rc *apiResponseCollector) Write(p []byte) (int, error) {
	return rc.stringbuilder.Write(p)
}

func (rc *apiResponseCollector) WriteHeader(statusCode int) {
	rc.statusCode = statusCode
}

func (rc *apiResponseCollector) ToAPIGatewayProxyResponse() events.APIGatewayProxyResponse {
	response := events.APIGatewayProxyResponse{
		StatusCode: rc.statusCode,
		Body:       rc.stringbuilder.String(),
	}

	if rc.useMultiValueHeader {
		response.MultiValueHeaders = rc.header
	} else {
		response.Headers = make(map[string]string)
		for k, v := range rc.header {
			response.Headers[k] = strings.Join(v, ",")
		}
	}

	return response
}
