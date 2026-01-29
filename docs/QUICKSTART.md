# Quick Start Guide

## 前置要求

- Node.js 18+
- Go 1.21+
- AWS CLI configured
- SAM CLI installed
- PostgreSQL database (Supabase)
- DeepSeek API Key

## 数据库设置

1. 在 Supabase 创建数据库
2. 运行 schema:
```bash
psql $DATABASE_URL < database/schema.sql
```

## 后端部署

1. 安装依赖:
```bash
cd lambda
make deps
```

2. 配置参数:
```bash
cp samconfig.toml.example samconfig.toml
# 编辑 samconfig.toml，填入实际参数
```

3. 部署:
```bash
make deploy
```

4. 记录 API Gateway URL

## 前端开发

1. 安装依赖:
```bash
cd frontend
npm install
```

2. 配置环境变量:
```bash
cp .env.example .env
# 编辑 .env，填入 API URL
```

3. 启动开发服务器:
```bash
npm run dev
```

4. 访问 http://localhost:5173

## 前端部署到 Amplify

1. 在 AWS Amplify Console 创建新应用
2. 连接 Git 仓库
3. 设置构建配置（使用 amplify.yml）
4. 添加环境变量:
   - `VITE_API_BASE_URL`: Lambda API Gateway URL
   - `VITE_DID_LOGIN_API_URL`: DID Login API URL
   - `VITE_TASK_UI_URL`: Task UI URL
5. 部署

## 环境变量清单

### Lambda (samconfig.toml)
```
DatabaseURL=postgresql://...
DeepSeekAPIKey=sk-xxx
DeepSeekModel=deepseek-chat
DeepSeekMaxTokens=2000
JWTSecret=your-jwt-secret
TaskUIAPIURL=https://...
```

### Frontend (.env)
```
VITE_API_BASE_URL=https://xxx.execute-api.us-east-1.amazonaws.com/Prod
VITE_DID_LOGIN_API_URL=https://i149gvmuh8.execute-api.us-east-1.amazonaws.com/prod
VITE_TASK_UI_URL=https://main.dwxknd52rdeie.amplifyapp.com
```

## 测试

### 测试对话功能
```bash
curl -X POST https://your-api-url/chat \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [
      {"role": "user", "content": "我想做跨境电商"}
    ],
    "project_id": "your-project-id"
  }'
```

### 测试保存报告
```bash
curl -X POST https://your-api-url/save-report \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "your-project-id",
    "business_goal": "跨境电商",
    "recommendations": {...}
  }'
```

## 常见问题

### 1. CORS 错误
确保 Lambda 函数返回正确的 CORS headers

### 2. JWT 验证失败
检查 JWT_SECRET 是否与 DID Login 系统一致

### 3. DeepSeek API 超时
增加 Lambda timeout 或减少 max_tokens

### 4. 数据库连接失败
检查 DATABASE_URL 格式和网络访问权限

## 下一步

- 查看 [API 文档](./API.md)
- 查看 [提示词设计](./PROMPTS.md)
- 查看 [需求文档](./REQUIREMENTS.md)
