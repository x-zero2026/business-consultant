package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/x-zero/business-consultant/pkg/db"
	"github.com/x-zero/business-consultant/pkg/response"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Headers": "Content-Type,Authorization",
				"Access-Control-Allow-Methods": "GET,OPTIONS",
			},
		}, nil
	}

	draftID := request.PathParameters["id"]
	if draftID == "" {
		return response.Error(400, "Draft ID is required")
	}

	if err := db.InitDB(); err != nil {
		return response.Error(500, fmt.Sprintf("Database error: %v", err))
	}

	pool := db.GetPool()

	var taskData []byte
	var used bool
	var expiresAt interface{}

	err := pool.QueryRow(ctx, `
		SELECT task_data, used, expires_at
		FROM task_drafts
		WHERE draft_id = $1 AND expires_at > NOW()
	`, draftID).Scan(&taskData, &used, &expiresAt)

	if err != nil {
		return response.Error(404, "Draft not found or expired")
	}

	if used {
		return response.Error(410, "Draft already used")
	}

	// Mark as used
	_, err = pool.Exec(ctx, `
		UPDATE task_drafts SET used = true WHERE draft_id = $1
	`, draftID)

	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to mark draft as used: %v\n", err)
	}

	// Parse task data
	var taskDataMap map[string]interface{}
	if err := json.Unmarshal(taskData, &taskDataMap); err != nil {
		return response.Error(500, "Failed to parse task data")
	}

	return response.Success(taskDataMap)
}

func main() {
	lambda.Start(handler)
}
