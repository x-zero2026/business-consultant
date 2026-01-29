package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/x-zero/business-consultant/pkg/auth"
	"github.com/x-zero/business-consultant/pkg/deepseek"
	"github.com/x-zero/business-consultant/pkg/response"
)

type ChatRequest struct {
	Messages  []deepseek.Message `json:"messages"`
	ProjectID string             `json:"project_id"`
	Stream    bool               `json:"stream"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Validate JWT
	authHeader := request.Headers["Authorization"]
	if authHeader == "" {
		authHeader = request.Headers["authorization"]
	}
	claims, err := auth.ValidateToken(authHeader)
	if err != nil {
		return response.ErrorNoCORS(401, fmt.Sprintf("Invalid token: %v", err))
	}

	// Parse request body
	var req ChatRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return response.ErrorNoCORS(400, "Invalid request body")
	}

	if len(req.Messages) == 0 {
		return response.ErrorNoCORS(400, "Messages cannot be empty")
	}

	if req.ProjectID == "" {
		return response.ErrorNoCORS(400, "Project ID is required")
	}

	client := deepseek.NewClient()

	// Handle streaming request
	if req.Stream {
		var accumulated strings.Builder
		var chunks []string

		err := client.ChatStream(req.Messages, func(chunk string) error {
			accumulated.WriteString(chunk)
			chunks = append(chunks, chunk)
			return nil
		})

		if err != nil {
			return response.ErrorNoCORS(500, fmt.Sprintf("AI error: %v", err))
		}

		fullResponse := accumulated.String()

		// Try to parse as JSON
		var aiData map[string]interface{}
		if err := json.Unmarshal([]byte(fullResponse), &aiData); err != nil {
			// If not valid JSON, return as plain text
			return response.SuccessNoCORS(map[string]interface{}{
				"stage":    "error",
				"message":  "AI返回格式异常，请重试",
				"raw":      fullResponse,
				"user_did": claims.DID,
			})
		}

		aiData["user_did"] = claims.DID
		return response.SuccessNoCORS(aiData)
	}

	// Handle non-streaming request (original behavior)
	aiResponse, err := client.Chat(req.Messages)
	if err != nil {
		return response.ErrorNoCORS(500, fmt.Sprintf("AI error: %v", err))
	}

	// Parse AI response
	var aiData map[string]interface{}
	if err := json.Unmarshal([]byte(aiResponse), &aiData); err != nil {
		// If JSON parsing fails, return as plain text with error stage
		return response.SuccessNoCORS(map[string]interface{}{
			"stage":    "error",
			"message":  "AI返回格式异常，请重试",
			"raw":      aiResponse,
			"user_did": claims.DID,
		})
	}

	// Add user info to response
	aiData["user_did"] = claims.DID

	return response.SuccessNoCORS(aiData)
}

func main() {
	lambda.Start(handler)
}
