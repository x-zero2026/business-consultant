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

	// Get project_id from query params
	projectID := request.QueryStringParameters["project_id"]
	if projectID == "" {
		return response.Error(400, "Project ID is required")
	}

	// Initialize database
	if err := db.InitDB(); err != nil {
		return response.Error(500, fmt.Sprintf("Database error: %v", err))
	}

	pool := db.GetPool()

	// Query reports
	rows, err := pool.Query(ctx, `
		SELECT report_id, project_id, business_goal, recommendations, created_at, updated_at
		FROM business_reports
		WHERE user_did = $1 AND project_id = $2
		ORDER BY created_at DESC
	`, claims.DID, projectID)

	if err != nil {
		return response.Error(500, fmt.Sprintf("Failed to query reports: %v", err))
	}
	defer rows.Close()

	// Collect reports
	reports := []map[string]interface{}{}
	for rows.Next() {
		var reportID, projectIDVal, businessGoal string
		var recommendations []byte
		var createdAt, updatedAt interface{}

		err := rows.Scan(&reportID, &projectIDVal, &businessGoal, &recommendations, &createdAt, &updatedAt)
		if err != nil {
			return response.Error(500, fmt.Sprintf("Failed to scan report: %v", err))
		}

		// Parse recommendations JSON
		var recsMap map[string]interface{}
		if err := json.Unmarshal(recommendations, &recsMap); err != nil {
			recsMap = nil
		}

		reports = append(reports, map[string]interface{}{
			"report_id":       reportID,
			"project_id":      projectIDVal,
			"business_goal":   businessGoal,
			"recommendations": recsMap,
			"created_at":      createdAt,
			"updated_at":      updatedAt,
		})
	}

	return response.Success(reports)
}

func main() {
	lambda.Start(handler)
}
