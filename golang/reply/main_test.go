package main_test

import (
	"testing"

	main "example.com/main"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {

	tests := []struct {
		request events.APIGatewayProxyRequest
		err     error
	}{
		{
			// Test that the handler responds with the correct response
			// when a valid name is provided in the HTTP body
			request: events.APIGatewayProxyRequest{HTTPMethod: "POST"},
			expect:  "POST",
			err:     nil,
		},
	}

	for _, test := range tests {
		response, err := main.HandleRequest(nil, test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.Body)
	}

}
