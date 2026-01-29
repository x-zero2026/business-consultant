package response

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// Success returns a successful API Gateway response with CORS headers
func Success(data interface{}) (events.APIGatewayProxyResponse, error) {
	body, err := json.Marshal(map[string]interface{}{
		"success": true,
		"data":    data,
	})
	if err != nil {
		return Error(500, "Failed to marshal response")
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "Content-Type,Authorization",
			"Access-Control-Allow-Methods": "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		},
		Body: string(body),
	}, nil
}

// SuccessNoCORS returns a successful response without CORS headers (for Function URLs)
func SuccessNoCORS(data interface{}) (events.APIGatewayProxyResponse, error) {
	body, err := json.Marshal(map[string]interface{}{
		"success": true,
		"data":    data,
	})
	if err != nil {
		return ErrorNoCORS(500, "Failed to marshal response")
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}, nil
}

// Error returns an error API Gateway response with CORS headers
func Error(statusCode int, message string) (events.APIGatewayProxyResponse, error) {
	body, _ := json.Marshal(map[string]interface{}{
		"success": false,
		"error":   message,
	})

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "Content-Type,Authorization",
			"Access-Control-Allow-Methods": "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		},
		Body: string(body),
	}, nil
}

// ErrorNoCORS returns an error response without CORS headers (for Function URLs)
func ErrorNoCORS(statusCode int, message string) (events.APIGatewayProxyResponse, error) {
	body, _ := json.Marshal(map[string]interface{}{
		"success": false,
		"error":   message,
	})

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}, nil
}
