package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/x-zero/business-consultant/pkg/auth"
	"github.com/x-zero/business-consultant/pkg/db"
	"github.com/x-zero/business-consultant/pkg/response"
)

type UpdateItemRequest struct {
	Status *string `json:"status"`
	TaskID *string `json:"task_id"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Headers": "Content-Type,Authorization",
				"Access-Control-Allow-Methods": "PATCH,OPTIONS",
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

	reportID := request.PathParameters["id"]
	itemID := request.PathParameters["item_id"]
	if reportID == "" || itemID == "" {
		return response.Error(400, "Report ID and Item ID are required")
	}

	var req UpdateItemRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return response.Error(400, "Invalid request body")
	}

	if err := db.InitDB(); err != nil {
		return response.Error(500, fmt.Sprintf("Database error: %v", err))
	}

	pool := db.GetPool()

	// Get current report
	var recommendations []byte
	var userDID string
	err = pool.QueryRow(ctx, `
		SELECT user_did, recommendations FROM business_reports WHERE report_id = $1
	`, reportID).Scan(&userDID, &recommendations)

	if err != nil {
		return response.Error(404, "Report not found")
	}

	if userDID != claims.DID {
		return response.Error(403, "Access denied")
	}

	// Parse recommendations
	var recsMap map[string]interface{}
	if err := json.Unmarshal(recommendations, &recsMap); err != nil {
		return response.Error(500, "Failed to parse recommendations")
	}

	// Update item status (simplified - just store in a status map)
	// In production, you'd want a more structured approach
	if recsMap["item_statuses"] == nil {
		recsMap["item_statuses"] = make(map[string]interface{})
	}
	
	itemStatuses := recsMap["item_statuses"].(map[string]interface{})
	itemStatuses[itemID] = map[string]interface{}{
		"status":  req.Status,
		"task_id": req.TaskID,
	}

	// Save back
	updatedRecs, err := json.Marshal(recsMap)
	if err != nil {
		return response.Error(500, "Failed to marshal recommendations")
	}

	_, err = pool.Exec(ctx, `
		UPDATE business_reports
		SET recommendations = $1, updated_at = NOW()
		WHERE report_id = $2
	`, updatedRecs, reportID)

	if err != nil {
		return response.Error(500, fmt.Sprintf("Failed to update report: %v", err))
	}

	return response.Success(map[string]interface{}{
		"message": "Item updated successfully",
	})
}

func main() {
	lambda.Start(handler)
}
