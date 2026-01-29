# Business Consultant - 快速参考

## 🚀 一键注册

```bash
cd business-consultant/scripts
export JWT_TOKEN="your-jwt-token"
./register-app.sh
```

## 📋 常用命令

### 注册应用
```bash
# 本地环境
export JWT_TOKEN="your-token"
export APP_URL="http://localhost:5173"
./register-app.sh

# 生产环境
export JWT_TOKEN="your-token"
export APP_URL="https://your-amplify-url.amplifyapp.com"
./register-app.sh
```

### 检查应用
```bash
export JWT_TOKEN="your-token"
./check-app.sh
```

### 更新应用
```bash
# 重新运行注册脚本即可
./register-app.sh
# 提示时输入 'y' 确认更新
```

## 🔑 获取 JWT Token

1. 访问: https://main.d2fozf421c6ftf.amplifyapp.com
2. 登录账号
3. 按 F12 打开开发者工具
4. Application → Local Storage → 复制 `token` 值

## ✅ 验证成功

### 方法 1: 使用脚本
```bash
./check-app.sh
```

### 方法 2: 使用 API
```bash
curl -X GET "https://i149gvmuh8.execute-api.us-east-1.amazonaws.com/prod/api/apps" \
  -H "Authorization: Bearer $JWT_TOKEN" | jq '.data[] | select(.app_name == "商业顾问")'
```

### 方法 3: 在 UI 中
1. 访问 DID Login Dashboard
2. 选择任意项目
3. 查看应用列表
4. 应该看到 "🤖 商业顾问"

## 🐛 快速排查

| 问题 | 解决方案 |
|------|----------|
| 401 Unauthorized | 重新获取 token |
| App already exists | 输入 'y' 更新 |
| 403 Forbidden | 检查管理员权限 |
| Connection refused | 检查 API 可访问性 |
| 应用不显示 | 刷新页面，检查 is_global |

## 📊 应用信息

| 字段 | 值 |
|------|-----|
| 名称 | 商业顾问 |
| Emoji | 🤖 |
| 描述 | AI驱动的一人公司商业咨询服务 |
| is_global | true |
| 默认 URL | http://localhost:5173 |

## 🔗 快速链接

- **DID Login**: https://main.d2fozf421c6ftf.amplifyapp.com
- **API**: https://i149gvmuh8.execute-api.us-east-1.amazonaws.com/prod
- **文档**: [README.md](./README.md)
- **示例**: [EXAMPLE.md](./EXAMPLE.md)

## 💡 提示

- Token 有效期有限，过期需重新获取
- `is_global=true` 使应用对所有用户可见
- 部署到生产环境后记得更新 URL
- 使用 `check-app.sh` 验证注册状态

---

**需要详细文档？** 查看 [README.md](./README.md) 和 [EXAMPLE.md](./EXAMPLE.md)
