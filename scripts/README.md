# Business Consultant - App Registration Scripts

这个目录包含用于在 DID Login 系统中注册 Business Consultant 应用的脚本。

## 📋 文件说明

- **register-app.sh** - Shell 脚本，通过 API 注册应用（推荐）
- **register-app.sql** - SQL 脚本，直接在数据库中注册应用（备用方案）
- **README.md** - 本文档

## 🚀 快速开始

### 方法 1: 使用 Shell 脚本（推荐）

#### 1. 获取 JWT Token

访问 DID Login 系统并登录：
```
https://main.d2fozf421c6ftf.amplifyapp.com
```

打开浏览器开发者工具（F12）：
1. 进入 Application → Local Storage
2. 找到 `token` 字段
3. 复制 token 值

#### 2. 运行注册脚本

```bash
cd business-consultant/scripts

# 方式 A: 交互式输入 token
./register-app.sh

# 方式 B: 通过环境变量传递 token
export JWT_TOKEN="your-jwt-token-here"
./register-app.sh

# 方式 C: 自定义 APP URL（用于生产环境）
export JWT_TOKEN="your-jwt-token-here"
export APP_URL="https://your-amplify-url.amplifyapp.com"
./register-app.sh
```

#### 3. 验证注册成功

脚本成功后，你会看到：
```
✅ Success! App registered successfully.

Next steps:
1. Go to DID Login Dashboard: https://main.d2fozf421c6ftf.amplifyapp.com
2. Select any project
3. You should see '商业顾问 🤖' in the apps list (visible to all users)
4. Click on it to open: http://localhost:5173
```

### 方法 2: 使用 SQL 脚本（备用）

如果 API 方法不工作，可以直接在 Supabase 中执行 SQL：

#### 1. 获取你的 DID

```sql
SELECT did, username, email 
FROM users 
WHERE username = 'your-username';
```

#### 2. 编辑 SQL 脚本

打开 `register-app.sql`，替换：
```sql
created_by_did = 'your-did-here'  -- 替换为你的实际 DID
```

#### 3. 执行 SQL

在 Supabase SQL Editor 中：
1. 复制 `register-app.sql` 的内容
2. 粘贴到 SQL Editor
3. 点击 Run

#### 4. 验证

```sql
SELECT 
  app_id,
  app_name,
  emoji,
  url,
  is_global
FROM apps
WHERE app_name = '商业顾问';
```

## 🔧 配置选项

### 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `JWT_TOKEN` | (必需) | DID Login 系统的 JWT token |
| `API_URL` | `https://i149gvmuh8.execute-api.us-east-1.amazonaws.com/prod` | DID Login API 地址 |
| `APP_URL` | `http://localhost:5173` | Business Consultant 前端地址 |
| `APP_EMOJI` | `🤖` | 应用图标 emoji |
| `IS_GLOBAL` | `true` | 是否全局可见 |

### 应用配置

```json
{
  "app_name": "商业顾问",
  "app_description": "AI驱动的一人公司商业咨询服务，帮助规划资源配置和预算",
  "emoji": "🤖",
  "url": "http://localhost:5173",
  "is_global": true
}
```

**重要**: `is_global: true` 表示这个应用对所有用户和所有项目可见。

## 📊 注册流程

```
1. 用户登录 DID Login 系统
   ↓
2. 获取 JWT Token
   ↓
3. 运行注册脚本
   ↓
4. 脚本调用 POST /api/apps
   ↓
5. 创建应用记录（is_global=true）
   ↓
6. 应用在所有项目中可见
   ↓
7. 用户可以从任何项目访问商业顾问
```

## 🔍 验证步骤

### 1. 检查应用是否创建

```bash
curl -X GET "https://i149gvmuh8.execute-api.us-east-1.amazonaws.com/prod/api/apps" \
  -H "Authorization: Bearer $JWT_TOKEN" | jq '.data[] | select(.app_name == "商业顾问")'
```

### 2. 检查全局状态

```sql
SELECT app_name, is_global, url 
FROM apps 
WHERE app_name = '商业顾问';
```

应该返回：
```
 app_name | is_global |           url            
----------+-----------+--------------------------
 商业顾问  |     t     | http://localhost:5173
```

### 3. 在 UI 中验证

1. 访问 DID Login Dashboard
2. 选择任意项目
3. 在应用列表中应该看到 "🤖 商业顾问"
4. 点击应用图标，应该跳转到 Business Consultant

## 🐛 故障排查

### 问题 1: JWT Token 无效

**错误**: `401 Unauthorized` 或 `Invalid token`

**解决**:
1. 确认 token 没有过期
2. 重新登录获取新 token
3. 检查 token 格式（应该是 `eyJ...` 开头）

### 问题 2: 应用已存在

**错误**: `App already exists`

**解决**:
1. 脚本会提示是否更新
2. 输入 `y` 更新现有应用
3. 或者手动删除后重新创建

### 问题 3: 权限不足

**错误**: `403 Forbidden` 或 `Access denied`

**解决**:
1. 确认你有管理员权限
2. 检查是否登录了正确的账号
3. 联系系统管理员授予权限

### 问题 4: API 无法访问

**错误**: `Connection refused` 或 `Timeout`

**解决**:
```bash
# 测试 API 是否可访问
curl https://i149gvmuh8.execute-api.us-east-1.amazonaws.com/prod/health

# 检查网络连接
ping i149gvmuh8.execute-api.us-east-1.amazonaws.com
```

### 问题 5: 应用不显示

**症状**: 应用创建成功，但在 UI 中看不到

**检查**:
1. 确认 `is_global = true`
2. 刷新浏览器页面
3. 清除浏览器缓存
4. 检查数据库中的记录

```sql
-- 检查应用状态
SELECT * FROM apps WHERE app_name = '商业顾问';

-- 检查是否有项目关联（全局应用不需要）
SELECT * FROM app_projects WHERE app_id IN (
  SELECT app_id FROM apps WHERE app_name = '商业顾问'
);
```

## 📝 更新应用信息

### 更新 URL（部署到生产环境后）

```bash
# 方法 1: 重新运行脚本
export JWT_TOKEN="your-token"
export APP_URL="https://your-production-url.amplifyapp.com"
./register-app.sh

# 方法 2: 使用 SQL
UPDATE apps
SET 
  url = 'https://your-production-url.amplifyapp.com',
  updated_at = NOW()
WHERE app_name = '商业顾问';
```

### 更新描述或 Emoji

```bash
# 编辑脚本中的配置
APP_DESCRIPTION="新的描述"
APP_EMOJI="💼"

# 重新运行
./register-app.sh
```

## 🎯 生产环境部署

### 1. 部署前端到 Amplify

```bash
cd business-consultant/frontend
npm run build
# 通过 Amplify Console 部署
```

### 2. 获取 Amplify URL

部署成功后，记录 Amplify URL，例如：
```
https://main.d1234567890abc.amplifyapp.com
```

### 3. 更新应用 URL

```bash
export JWT_TOKEN="your-token"
export APP_URL="https://main.d1234567890abc.amplifyapp.com"
./register-app.sh
```

### 4. 验证

访问 DID Login Dashboard，点击商业顾问应用，应该跳转到生产环境 URL。

## 📚 相关文档

- [EXAMPLE.md](./EXAMPLE.md) - 详细使用示例和场景
- [Business Consultant README](../README.md) - 项目主文档
- [部署指南](../DEPLOYMENT-CHECKLIST.md) - 完整部署清单
- [快速开始](../docs/QUICKSTART.md) - 快速开始指南
- [应用注册完成报告](../APP-REGISTRATION-COMPLETE.md) - 注册脚本说明
- [DID Login API 文档](../../docs/did-login/API-REFERENCE.md) - API 参考文档

## 💡 提示

1. **全局应用**: `is_global=true` 使应用对所有用户可见，无需为每个项目单独添加
2. **Token 安全**: 不要将 JWT token 提交到 Git 仓库
3. **URL 更新**: 部署到生产环境后记得更新应用 URL
4. **测试**: 在生产环境注册前，先在本地测试脚本

## 🔗 相关链接

- **DID Login Dashboard**: https://main.d2fozf421c6ftf.amplifyapp.com
- **DID Login API**: https://i149gvmuh8.execute-api.us-east-1.amazonaws.com/prod
- **Business Consultant (Local)**: http://localhost:5173
- **Supabase Dashboard**: https://supabase.com/dashboard

---

**需要帮助？** 查看 [故障排查](#-故障排查) 部分或联系开发团队。
