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

type IdentifyTagsRequest struct {
	TaskDescription string `json:"task_description"`
}

type IdentifyTagsResponse struct {
	ProfessionTags []string `json:"profession_tags"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("=== IdentifyProfessionTags Handler Started ===")

	// Handle CORS preflight
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

	// Extract and validate token
	authHeader := request.Headers["Authorization"]
	if authHeader == "" {
		authHeader = request.Headers["authorization"]
	}

	_, err := auth.ValidateToken(authHeader)
	if err != nil {
		fmt.Printf("Token validation error: %v\n", err)
		return response.Error(401, "Invalid or expired token")
	}

	// Parse request body
	var req IdentifyTagsRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		fmt.Printf("JSON parse error: %v\n", err)
		return response.Error(400, "Invalid request body")
	}

	if req.TaskDescription == "" {
		return response.Error(400, "task_description is required")
	}

	fmt.Printf("Task description: %s\n", req.TaskDescription)

	// Call DeepSeek to identify tags
	tags, err := identifyTags(ctx, req.TaskDescription)
	if err != nil {
		fmt.Printf("DeepSeek API error: %v\n", err)
		// Return empty array on error
		return response.Success(IdentifyTagsResponse{
			ProfessionTags: []string{},
		})
	}

	fmt.Printf("Identified tags: %v\n", tags)

	return response.Success(IdentifyTagsResponse{
		ProfessionTags: tags,
	})
}

func identifyTags(ctx context.Context, taskDescription string) ([]string, error) {
	messages := []deepseek.Message{
		{
			Role: "system",
			Content: `你是一个职业标签识别专家。根据任务描述，识别出最相关的职业标签。

预定义标签列表（优先使用）：
- frontend-developer, backend-developer, fullstack-developer, mobile-developer
- devops-engineer, data-engineer, ml-engineer, qa-engineer
- ui-designer, ux-designer, product-designer, graphic-designer
- product-manager, project-manager, scrum-master
- business-analyst, entrepreneur, consultant
- researcher, writer, marketer

要求：
1. 优先从预定义标签中选择
2. 如果预定义标签不够准确，可以添加自定义标签（如：客服专员、运营经理、销售总监）
3. 最多返回5个标签
4. 只返回JSON格式，不要有其他文字：{"profession_tags": ["tag1", "tag2"]}
5. 标签要准确反映任务所需的职业技能`,
		},
		{
			Role:    "user",
			Content: fmt.Sprintf("任务描述：\n%s", taskDescription),
		},
	}

	client := deepseek.NewClient()
	content, err := client.Chat(messages)
	if err != nil {
		return nil, err
	}

	// Parse JSON response
	content = strings.TrimSpace(content)
	
	// Try to extract JSON if wrapped in markdown code blocks
	if strings.HasPrefix(content, "```json") {
		content = strings.TrimPrefix(content, "```json")
		content = strings.TrimSuffix(content, "```")
		content = strings.TrimSpace(content)
	} else if strings.HasPrefix(content, "```") {
		content = strings.TrimPrefix(content, "```")
		content = strings.TrimSuffix(content, "```")
		content = strings.TrimSpace(content)
	}

	var result struct {
		ProfessionTags []string `json:"profession_tags"`
	}

	if err := json.Unmarshal([]byte(content), &result); err != nil {
		fmt.Printf("Failed to parse DeepSeek response: %s\n", content)
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Limit to 5 tags
	if len(result.ProfessionTags) > 5 {
		result.ProfessionTags = result.ProfessionTags[:5]
	}

	return result.ProfessionTags, nil
}

func main() {
	lambda.Start(handler)
}
