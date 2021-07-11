package ginlambda

import (
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type albResponseCollector struct {
	header              map[string][]string
	stringbuilder       strings.Builder
	statusCode          int
	useMultiValueHeader bool
}

func newALBResponseCollector(useMultiValueHeader bool) *albResponseCollector {
	return &albResponseCollector{
		header:              make(map[string][]string),
		useMultiValueHeader: useMultiValueHeader,
	}
}

func (rc *albResponseCollector) Header() http.Header {
	return rc.header
}

func (rc *albResponseCollector) Write(p []byte) (int, error) {
	return rc.stringbuilder.Write(p)
}

func (rc *albResponseCollector) WriteHeader(statusCode int) {
	rc.statusCode = statusCode
}

func (rc *albResponseCollector) ToAPIGatewayProxyResponse() events.APIGatewayProxyResponse {
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
