# Business Consultant - 一人公司智能商业顾问

为一人公司创业者提供AI驱动的商业咨询服务，帮助规划资源分配、推荐AI工作流和真人岗位。

## 功能特性

- 🤖 **AI对话咨询**: 基于 DeepSeek API 的多轮对话
- 📊 **智能推荐**: 自动生成 AI 工作流和真人岗位建议
- 💰 **预算规划**: 分阶段的月度预算明细
- 📝 **报告管理**: 保存和查看历史咨询报告
- 🔗 **任务集成**: 一键发布任务到 Task UI
- 🔐 **DID 登录**: 集成 X-Zero 统一身份认证

## 技术栈

### 前端
- React 18 + Vite
- React Router v6
- Axios
- jsPDF (导出功能)

### 后端
- Go 1.21
- AWS Lambda + API Gateway
- PostgreSQL (Supabase)
- DeepSeek API

## 项目结构

```
business-consultant/
├── frontend/           # React 前端
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── api/
│   │   └── utils/
│   └── package.json
├── lambda/            # Go Lambda 函数
│   ├── cmd/
│   ├── pkg/
│   └── template.yaml
├── database/          # 数据库脚本
│   └── schema.sql
└── docs/             # 文档
    ├── REQUIREMENTS.md
    ├── API.md
    └── PROMPTS.md
```

## 快速开始

### 前端开发

```bash
cd frontend
npm install
npm run dev
```

### 后端开发

```bash
cd lambda
sam build
sam local start-api
```

## 环境变量

### Frontend (.env)
```
VITE_API_BASE_URL=https://api.business-consultant.com
VITE_DID_LOGIN_API_URL=https://did-login.com/api
VITE_TASK_UI_URL=https://task-ui.com
```

### Lambda (SAM template)
```
DATABASE_URL=postgresql://...
DEEPSEEK_API_KEY=sk-xxx
DEEPSEEK_MODEL=deepseek-chat
JWT_SECRET=xxx
TASK_UI_API_URL=https://task-ui.com/api
```

## 部署

### 1. 后端部署到 AWS
```bash
cd lambda
sam build
sam deploy --guided
```

### 2. 前端部署到 Amplify
```bash
cd frontend
npm run build
# 通过 Amplify Console 部署
```

### 3. 注册应用到 DID Login 系统
```bash
cd scripts

# 获取 JWT token 后运行
export JWT_TOKEN="your-jwt-token"
export APP_URL="https://your-amplify-url.amplifyapp.com"
./register-app.sh
```

详细说明请查看 [scripts/README.md](./scripts/README.md)

## 开发进度

- [x] 项目初始化
- [x] 对话界面
- [x] AI 推荐生成
- [x] 报告管理
- [x] 任务发布集成
- [ ] 部署上线

## License

MIT
