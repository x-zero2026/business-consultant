package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/x-zero/business-consultant/pkg/auth"
	"github.com/x-zero/business-consultant/pkg/db"
	"github.com/x-zero/business-consultant/pkg/response"
)

type CreateDraftRequest struct {
	DraftID   string                 `json:"draft_id"`
	TaskData  map[string]interface{} `json:"task_data"`
	ExpiresAt string                 `json:"expires_at"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Headers": "Content-Type,Authorization",
				"Access-Control-Allow-Methods": "POST,OPTIONS",
			},
		}, nil
	}

	authHeader := request.Headers["Authorization"]
	if authHeader == "" {
		authHeader = request.Headers["authorization"]
	}
	claims, err := auth.ValidateToken(authHeader)
	if err != nil {
		return response.Error(401, fmt.Sprintf("Invalid token: %v", err))
	}

	var req CreateDraftRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return response.Error(400, "Invalid request body")
	}

	if req.DraftID == "" {
		return response.Error(400, "Draft ID is required")
	}

	// Parse expires_at
	expiresAt, err := time.Parse(time.RFC3339, req.ExpiresAt)
	if err != nil {
		// Default to 1 hour from now
		expiresAt = time.Now().Add(time.Hour)
	}

	if err := db.InitDB(); err != nil {
		return response.Error(500, fmt.Sprintf("Database error: %v", err))
	}

	pool := db.GetPool()

	taskDataJSON, err := json.Marshal(req.TaskData)
	if err != nil {
		return response.Error(500, "Failed to marshal task data")
	}

	_, err = pool.Exec(ctx, `
		INSERT INTO task_drafts (draft_id, user_did, task_data, expires_at)
		VALUES ($1, $2, $3, $4)
	`, req.DraftID, claims.DID, taskDataJSON, expiresAt)

	if err != nil {
		return response.Error(500, fmt.Sprintf("Failed to create draft: %v", err))
	}

	return response.Success(map[string]interface{}{
		"draft_id": req.DraftID,
		"message":  "Draft created successfully",
	})
}

func main() {
	lambda.Start(handler)
}
