package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/x-zero/business-consultant/pkg/auth"
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
				"Access-Control-Allow-Methods": "DELETE,OPTIONS",
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
	if reportID == "" {
		return response.Error(400, "Report ID is required")
	}

	if err := db.InitDB(); err != nil {
		return response.Error(500, fmt.Sprintf("Database error: %v", err))
	}

	pool := db.GetPool()

	result, err := pool.Exec(ctx, `
		DELETE FROM business_reports
		WHERE report_id = $1 AND user_did = $2
	`, reportID, claims.DID)

	if err != nil {
		return response.Error(500, fmt.Sprintf("Failed to delete report: %v", err))
	}

	if result.RowsAffected() == 0 {
		return response.Error(404, "Report not found or access denied")
	}

	return response.Success(map[string]interface{}{
		"message": "Report deleted successfully",
	})
}

func main() {
	lambda.Start(handler)
}
