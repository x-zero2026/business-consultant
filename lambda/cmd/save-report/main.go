package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"github.com/x-zero/business-consultant/pkg/auth"
	"github.com/x-zero/business-consultant/pkg/db"
	"github.com/x-zero/business-consultant/pkg/response"
)

type SaveReportRequest struct {
	ProjectID       string                 `json:"project_id"`
	BusinessGoal    string                 `json:"business_goal"`
	Recommendations map[string]interface{} `json:"recommendations"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Handle OPTIONS
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

	// Validate JWT
	authHeader := request.Headers["Authorization"]
	if authHeader == "" {
		authHeader = request.Headers["authorization"]
	}
	claims, err := auth.ValidateToken(authHeader)
	if err != nil {
		return response.Error(401, fmt.Sprintf("Invalid token: %v", err))
	}

	// Parse request body
	var req SaveReportRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return response.Error(400, "Invalid request body")
	}

	if req.ProjectID == "" {
		return response.Error(400, "Project ID is required")
	}

	if req.BusinessGoal == "" {
		return response.Error(400, "Business goal is required")
	}

	// Initialize database
	if err := db.InitDB(); err != nil {
		return response.Error(500, fmt.Sprintf("Database error: %v", err))
	}

	pool := db.GetPool()

	// Convert recommendations to JSON string
	recommendationsJSON, err := json.Marshal(req.Recommendations)
	if err != nil {
		return response.Error(500, "Failed to marshal recommendations")
	}

	// Insert report
	reportID := uuid.New().String()
	_, err = pool.Exec(ctx, `
		INSERT INTO business_reports (report_id, user_did, project_id, business_goal, recommendations)
		VALUES ($1, $2, $3, $4, $5::jsonb)
	`, reportID, claims.DID, req.ProjectID, req.BusinessGoal, string(recommendationsJSON))

	if err != nil {
		return response.Error(500, fmt.Sprintf("Failed to save report: %v", err))
	}

	return response.Success(map[string]interface{}{
		"report_id": reportID,
		"message":   "Report saved successfully",
	})
}

func main() {
	lambda.Start(handler)
}
