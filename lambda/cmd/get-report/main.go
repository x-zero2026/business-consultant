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

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Handle OPTIONS
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

	// Validate JWT
	authHeader := request.Headers["Authorization"]
	if authHeader == "" {
		authHeader = request.Headers["authorization"]
	}
	claims, err := auth.ValidateToken(authHeader)
	if err != nil {
		return response.Error(401, fmt.Sprintf("Invalid token: %v", err))
	}

	// Get report ID from path
	reportID := request.PathParameters["id"]
	if reportID == "" {
		return response.Error(400, "Report ID is required")
	}

	// Initialize database
	if err := db.InitDB(); err != nil {
		return response.Error(500, fmt.Sprintf("Database error: %v", err))
	}

	pool := db.GetPool()

	// Query report
	var projectID, businessGoal, userDID string
	var recommendations []byte
	var createdAt, updatedAt interface{}

	err = pool.QueryRow(ctx, `
		SELECT report_id, user_did, project_id, business_goal, recommendations, created_at, updated_at
		FROM business_reports
		WHERE report_id = $1
	`, reportID).Scan(&reportID, &userDID, &projectID, &businessGoal, &recommendations, &createdAt, &updatedAt)

	if err != nil {
		return response.Error(404, "Report not found")
	}

	// Verify ownership
	if userDID != claims.DID {
		return response.Error(403, "Access denied")
	}

	// Parse recommendations JSON
	var recsMap map[string]interface{}
	if err := json.Unmarshal(recommendations, &recsMap); err != nil {
		return response.Error(500, "Failed to parse recommendations")
	}

	return response.Success(map[string]interface{}{
		"report_id":       reportID,
		"user_did":        userDID,
		"project_id":      projectID,
		"business_goal":   businessGoal,
		"recommendations": recsMap,
		"created_at":      createdAt,
		"updated_at":      updatedAt,
	})
}

func main() {
	lambda.Start(handler)
}
