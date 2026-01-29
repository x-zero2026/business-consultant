package deepseek

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	DeepSeekAPIURL = "https://api.deepseek.com/v1/chat/completions"
)

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents the request to DeepSeek API
type ChatRequest struct {
	Model          string    `json:"model"`
	Messages       []Message `json:"messages"`
	Temperature    float64   `json:"temperature"`
	MaxTokens      int       `json:"max_tokens"`
	Stream         bool      `json:"stream"`
	ResponseFormat *struct {
		Type string `json:"type"`
	} `json:"response_format,omitempty"`
}

// ChatResponse represents the response from DeepSeek API
type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// Client represents a DeepSeek API client
type Client struct {
	APIKey     string
	Model      string
	MaxTokens  int
	HTTPClient *http.Client
}

// NewClient creates a new DeepSeek client
func NewClient() *Client {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	model := os.Getenv("DEEPSEEK_MODEL")
	if model == "" {
		model = "deepseek-chat"
	}

	maxTokens := 2000
	if mt := os.Getenv("DEEPSEEK_MAX_TOKENS"); mt != "" {
		fmt.Sscanf(mt, "%d", &maxTokens)
	}

	return &Client{
		APIKey:     apiKey,
		Model:      model,
		MaxTokens:  maxTokens,
		HTTPClient: &http.Client{},
	}
}

// Chat sends a chat request to DeepSeek API
func (c *Client) Chat(messages []Message) (string, error) {
	if c.APIKey == "" {
		return "", fmt.Errorf("DEEPSEEK_API_KEY not set")
	}

	// Add system prompt if not present
	if len(messages) == 0 || messages[0].Role != "system" {
		systemPrompt := getSystemPrompt()
		messages = append([]Message{{Role: "system", Content: systemPrompt}}, messages...)
	}

	// Prepare request
	reqBody := ChatRequest{
		Model:       c.Model,
		Messages:    messages,
		Temperature: 0.7,
		MaxTokens:   c.MaxTokens,
		Stream:      false,
		ResponseFormat: &struct {
			Type string `json:"type"`
		}{
			Type: "json_object",
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", DeepSeekAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	// Send request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	content := chatResp.Choices[0].Message.Content

	// Try to parse and fix JSON if needed
	content = fixJSON(content)

	return content, nil
}

// ChatStream sends a streaming chat request to DeepSeek API
func (c *Client) ChatStream(messages []Message, callback func(string) error) error {
	if c.APIKey == "" {
		return fmt.Errorf("DEEPSEEK_API_KEY not set")
	}

	// Add system prompt if not present
	if len(messages) == 0 || messages[0].Role != "system" {
		systemPrompt := getSystemPrompt()
		messages = append([]Message{{Role: "system", Content: systemPrompt}}, messages...)
	}

	// Prepare request (no response_format for streaming)
	reqBody := ChatRequest{
		Model:       c.Model,
		Messages:    messages,
		Temperature: 0.7,
		MaxTokens:   c.MaxTokens,
		Stream:      true,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", DeepSeekAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Accept", "text/event-stream")

	// Send request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Read streaming response
	reader := io.Reader(resp.Body)
	buf := make([]byte, 4096)
	var accumulated strings.Builder

	for {
		n, err := reader.Read(buf)
		if n > 0 {
			chunk := string(buf[:n])
			lines := strings.Split(chunk, "\n")

			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" || line == "data: [DONE]" {
					continue
				}

				if strings.HasPrefix(line, "data: ") {
					data := strings.TrimPrefix(line, "data: ")
					
					var streamResp struct {
						Choices []struct {
							Delta struct {
								Content string `json:"content"`
							} `json:"delta"`
						} `json:"choices"`
					}

					if err := json.Unmarshal([]byte(data), &streamResp); err == nil {
						if len(streamResp.Choices) > 0 {
							content := streamResp.Choices[0].Delta.Content
							if content != "" {
								accumulated.WriteString(content)
								if err := callback(content); err != nil {
									return err
								}
							}
						}
					}
				}
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read stream: %v", err)
		}
	}

	return nil
}

// fixJSON attempts to extract and fix JSON from the response
func fixJSON(text string) string {
	// Remove markdown code blocks if present
	text = strings.TrimSpace(text)
	
	// Remove ```json and ``` markers
	if strings.HasPrefix(text, "```json") {
		text = strings.TrimPrefix(text, "```json")
		text = strings.TrimSuffix(text, "```")
		text = strings.TrimSpace(text)
	} else if strings.HasPrefix(text, "```") {
		text = strings.TrimPrefix(text, "```")
		text = strings.TrimSuffix(text, "```")
		text = strings.TrimSpace(text)
	}

	return text
}

// getSystemPrompt returns the system prompt for the AI
func getSystemPrompt() string {
	return `你是一位专业的一人公司商业顾问，专门帮助创业者规划资源配置和预算。

**重要：你必须始终返回纯JSON格式的响应。**

## 对话流程

### 阶段1：询问商业目标
返回JSON格式：
{
  "stage": "initial",
  "message": "您好！我是您的一人公司商业顾问。请告诉我您的商业目标是什么？例如：跨境电商、SaaS产品、内容创作等。"
}

### 阶段2：询问细节
返回JSON格式：
{
  "stage": "questioning",
  "message": "了解了，您打算做[商业目标]。为了给您更精准的建议，我需要了解一些细节：",
  "questions": [
    "您的目标市场是哪些国家或地区？",
    "您的初始预算大约是多少？（单位：人民币）",
    "您是否有相关行业经验？"
  ],
  "progress": "2/3"
}

### 阶段3：提供推荐
返回JSON格式（必须包含所有字段）：
{
  "stage": "recommending",
  "message": "根据您的情况，我为您制定了以下方案：",
  "recommendations": {
    "business_goal": "跨境电商",
    "summary": "基于您的预算和背景，建议采用轻资产模式...",
    "ai_workflows": [
      {
        "name": "商品图片AI优化",
        "description": "使用AI自动优化商品图片，提升转化率",
        "input_requirements": "商品原图（JPG/PNG格式）",
        "output_requirements": "优化后的商品图（统一尺寸、背景、色调）",
        "estimated_cost": 100,
        "priority": "high"
      }
    ],
    "human_roles": [
      {
        "title": "客服专员（兼职）",
        "responsibilities": ["处理客户咨询", "订单跟进"],
        "requirements": ["良好的沟通能力", "熟悉电商平台"],
        "work_hours": "每天2-3小时，灵活安排",
        "monthly_budget": 2000,
        "priority": "high"
      }
    ],
    "phases": [
      {
        "phase_name": "启动期（第1-3个月）",
        "duration": "3个月",
        "monthly_budget": 3500,
        "budget_breakdown": {
          "AI工具": 100,
          "客服": 2000,
          "营销": 1000,
          "其他": 400
        }
      }
    ]
  }
}

## 重要规则

1. **必须严格返回纯JSON格式**
2. 不要添加Markdown标记或代码块符号
3. 预算单位为XZT（1 XZT ≈ 1 CNY）
4. 优先推荐AI自动化方案
5. 只讨论商业相关话题
6. **所有字段都必须填写，不能为空**：
   - ai_workflows必须包含：name, description, input_requirements, output_requirements, estimated_cost(数字), priority
   - human_roles必须包含：title, responsibilities, requirements, work_hours, monthly_budget(数字), priority
   - phases必须包含：phase_name, duration, monthly_budget(数字), budget_breakdown(对象)
7. **所有预算字段必须是数字，不能是文字描述**：
   - estimated_cost: 必须是数字（如100）
   - monthly_budget: 必须是数字（如3500）
   - budget_breakdown: 必须是对象，值必须是数字（如{"AI工具": 100, "客服": 2000}）

## 预算参考（XZT）
- AI工具: 20-200/月
- 兼职客服: 1500-3000/月
- 兼职开发: 3000-8000/月
- 营销推广: 500-2000/月`
}
