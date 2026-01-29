# Business Consultant - 注册示例

## 完整注册流程示例

### 场景 1: 本地开发环境注册

```bash
# 1. 进入脚本目录
cd business-consultant/scripts

# 2. 设置环境变量
export JWT_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
export APP_URL="http://localhost:5173"

# 3. 运行注册脚本
./register-app.sh
```

**输出示例**:
```
========================================
Business Consultant - App Registration
========================================

Configuration:
  API URL: https://i149gvmuh8.execute-api.us-east-1.amazonaws.com/prod
  App Name: 商业顾问
  App URL: http://localhost:5173
  Description: AI驱动的一人公司商业咨询服务，帮助规划资源配置和预算
  Emoji: 🤖
  Is Global: true

Checking if app already exists...
Creating app...

✅ Success! App registered successfully.

Response:
{
  "success": true,
  "data": {
    "app_id": "123e4567-e89b-12d3-a456-426614174000",
    "app_name": "商业顾问",
    "url": "http://localhost:5173",
    "is_global": true
  }
}

Next steps:
1. Go to DID Login Dashboard: https://main.d2fozf421c6ftf.amplifyapp.com
2. Select any project
3. You should see '商业顾问 🤖' in the apps list (visible to all users)
4. Click on it to open: http://localhost:5173
```

### 场景 2: 生产环境注册

```bash
# 1. 部署前端到 Amplify
cd business-consultant/frontend
npm run build
# 通过 Amplify Console 部署，获得 URL: https://main.d1234567890abc.amplifyapp.com

# 2. 注册应用
cd ../scripts
export JWT_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
export APP_URL="https://main.d1234567890abc.amplifyapp.com"
./register-app.sh
```

### 场景 3: 更新现有应用

```bash
# 如果应用已存在，脚本会提示是否更新
./register-app.sh

# 输出:
# ⚠ App '商业顾问' already exists (ID: 123e4567-e89b-12d3-a456-426614174000)
# 
# Do you want to update it? (y/n): y
# 
# Updating app...
# ✅ Success! App updated successfully.
```

### 场景 4: 使用 SQL 直接注册

```sql
-- 1. 获取你的 DID
SELECT did, username FROM users WHERE username = 'chilly';
-- 结果: did = '0x1234567890abcdef...'

-- 2. 插入应用
INSERT INTO apps (
  app_id,
  app_name,
  app_description,
  emoji,
  url,
  is_global,
  created_by_did,
  created_at,
  updated_at
) VALUES (
  gen_random_uuid(),
  '商业顾问',
  'AI驱动的一人公司商业咨询服务，帮助规划资源配置和预算',
  '🤖',
  'http://localhost:5173',
  true,
  '0x1234567890abcdef...',  -- 替换为你的 DID
  NOW(),
  NOW()
);

-- 3. 验证
SELECT app_name, is_global, url FROM apps WHERE app_name = '商业顾问';
```

## 验证注册成功

### 方法 1: 使用检查脚本

```bash
cd business-consultant/scripts
export JWT_TOKEN="your-token"
./check-app.sh
```

**成功输出**:
```
========================================
Business Consultant - Check Registration
========================================

Checking app registration...

✅ App is registered!

App Details:
  Name: 商业顾问 🤖
  ID: 123e4567-e89b-12d3-a456-426614174000
  URL: http://localhost:5173
  Description: AI驱动的一人公司商业咨询服务，帮助规划资源配置和预算
  Is Global: true
  Created: 2026-01-29T10:30:00.000Z

✓ App is set as global (visible to all users)

Access the app:
1. Go to: https://main.d2fozf421c6ftf.amplifyapp.com
2. Login and select any project
3. Click on '商业顾问 🤖' in the apps list
4. You will be redirected to: http://localhost:5173
```

### 方法 2: 使用 API

```bash
curl -X GET "https://i149gvmuh8.execute-api.us-east-1.amazonaws.com/prod/api/apps" \
  -H "Authorization: Bearer $JWT_TOKEN" | jq '.data[] | select(.app_name == "商业顾问")'
```

**成功输出**:
```json
{
  "app_id": "123e4567-e89b-12d3-a456-426614174000",
  "app_name": "商业顾问",
  "app_description": "AI驱动的一人公司商业咨询服务，帮助规划资源配置和预算",
  "emoji": "🤖",
  "url": "http://localhost:5173",
  "is_global": true,
  "created_by_did": "0x1234567890abcdef...",
  "created_at": "2026-01-29T10:30:00.000Z",
  "updated_at": "2026-01-29T10:30:00.000Z"
}
```

### 方法 3: 在 UI 中验证

1. 访问 DID Login Dashboard: https://main.d2fozf421c6ftf.amplifyapp.com
2. 登录你的账号
3. 选择任意项目
4. 在应用列表中应该看到 "🤖 商业顾问"
5. 点击应用图标
6. 应该跳转到 Business Consultant 页面

## 常见场景处理

### 场景 A: Token 过期

```bash
# 错误: 401 Unauthorized

# 解决:
# 1. 重新登录 DID Login
# 2. 获取新的 token
# 3. 重新运行脚本
export JWT_TOKEN="new-token-here"
./register-app.sh
```

### 场景 B: 应用名称冲突

```bash
# 错误: App already exists

# 解决方案 1: 更新现有应用
./register-app.sh
# 输入 'y' 确认更新

# 解决方案 2: 删除后重新创建
# 在 Supabase SQL Editor 中:
DELETE FROM apps WHERE app_name = '商业顾问';
# 然后重新运行脚本
./register-app.sh
```

### 场景 C: 需要修改应用信息

```bash
# 修改 URL
export JWT_TOKEN="your-token"
export APP_URL="https://new-url.amplifyapp.com"
./register-app.sh

# 修改 Emoji
# 编辑 register-app.sh，修改 APP_EMOJI 变量
APP_EMOJI="💼"
./register-app.sh
```

### 场景 D: 设置为非全局应用

```bash
# 编辑 register-app.sh
IS_GLOBAL="false"

# 运行脚本
./register-app.sh

# 然后需要为特定项目添加应用
# 在 Supabase SQL Editor 中:
INSERT INTO app_projects (app_id, project_id, added_at)
SELECT 
  (SELECT app_id FROM apps WHERE app_name = '商业顾问'),
  'your-project-id',
  NOW();
```

## 批量操作示例

### 为多个环境注册

```bash
#!/bin/bash

# 开发环境
export JWT_TOKEN="dev-token"
export APP_URL="http://localhost:5173"
./register-app.sh

# 测试环境
export JWT_TOKEN="test-token"
export APP_URL="https://test.amplifyapp.com"
./register-app.sh

# 生产环境
export JWT_TOKEN="prod-token"
export APP_URL="https://prod.amplifyapp.com"
./register-app.sh
```

### 检查所有全局应用

```sql
SELECT 
  app_name,
  emoji,
  url,
  is_global,
  created_at
FROM apps
WHERE is_global = true
ORDER BY created_at DESC;
```

**结果示例**:
```
   app_name   | emoji |              url              | is_global |       created_at        
--------------+-------+-------------------------------+-----------+-------------------------
 商业顾问      | 🤖    | http://localhost:5173         |     t     | 2026-01-29 10:30:00
 AI工作流中心  | 🤖    | http://localhost:5174         |     t     | 2026-01-28 15:20:00
 人才市场      | 💰    | http://localhost:5175         |     t     | 2026-01-27 09:10:00
```

## 调试技巧

### 启用详细输出

```bash
# 使用 bash -x 查看详细执行过程
bash -x register-app.sh
```

### 查看 API 响应

```bash
# 手动调用 API 查看详细响应
curl -v -X POST "https://i149gvmuh8.execute-api.us-east-1.amazonaws.com/prod/api/apps" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "app_name": "商业顾问",
    "url": "http://localhost:5173",
    "app_description": "AI驱动的一人公司商业咨询服务，帮助规划资源配置和预算",
    "emoji": "🤖",
    "is_global": true
  }'
```

### 检查数据库状态

```sql
-- 查看应用详情
SELECT * FROM apps WHERE app_name = '商业顾问';

-- 查看应用创建者
SELECT 
  a.app_name,
  a.is_global,
  u.username,
  u.did
FROM apps a
LEFT JOIN users u ON a.created_by_did = u.did
WHERE a.app_name = '商业顾问';

-- 查看应用访问日志（如果有）
SELECT * FROM app_access_logs 
WHERE app_id = (SELECT app_id FROM apps WHERE app_name = '商业顾问')
ORDER BY accessed_at DESC
LIMIT 10;
```

## 最佳实践

1. **使用环境变量**: 不要在脚本中硬编码 token
2. **验证注册**: 每次注册后运行 `check-app.sh` 验证
3. **记录 App ID**: 保存应用 ID 以便后续管理
4. **更新文档**: 部署到生产环境后更新 URL
5. **定期检查**: 确保应用状态正常

## 相关命令速查

```bash
# 注册应用
./register-app.sh

# 检查应用
./check-app.sh

# 获取所有应用
curl -X GET "$API_URL/api/apps" -H "Authorization: Bearer $JWT_TOKEN"

# 更新应用
curl -X PUT "$API_URL/api/apps/$APP_ID" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"url": "new-url"}'

# 设置为全局
curl -X POST "$API_URL/api/apps/$APP_ID/set-global" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"is_global": true}'

# 删除应用
curl -X DELETE "$API_URL/api/apps/$APP_ID" \
  -H "Authorization: Bearer $JWT_TOKEN"
```

---

**提示**: 所有示例中的 token 和 URL 都需要替换为实际值。
