# API Usage Guide with CommonResponse

本文档说明如何使用统一的 CommonResponse 格式的 API。

## Response 格式

所有 API 返回统一的响应格式：

```json
{
  "error_code": 0,
  "error_msg": "success",
  "data": {
    "result": {
      // 实际的响应数据
    }
  }
}
```

### 字段说明

- `error_code`: 错误码，0 表示成功，非 0 表示错误
- `error_msg`: 错误消息，成功时为 "success"
- `data`: 数据载体，使用灵活的结构化数据

### 错误码定义

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 400 | 参数错误 |
| 401 | 未认证 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 409 | 资源已存在 |
| 429 | 请求过多 |
| 500 | 内部服务器错误 |
| 501 | 功能未实现 |

## API 示例

### 1. 创建用户 (CreateUser)

**请求:**
```bash
curl -X POST http://localhost:8088/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "user": {
      "email": "alice@example.com",
      "display_name": "Alice Smith",
      "phone_number": "+1234567890"
    }
  }'
```

**成功响应:**
```json
{
  "error_code": 0,
  "error_msg": "success",
  "data": {
    "result": {
      "name": "users/1",
      "email": "alice@example.com",
      "display_name": "Alice Smith",
      "phone_number": "+1234567890",
      "create_time": "2025-12-17T10:00:00Z",
      "update_time": "2025-12-17T10:00:00Z",
      "is_active": true
    }
  }
}
```

**错误响应 (参数错误):**
```json
{
  "error_code": 400,
  "error_msg": "email is required",
  "data": null
}
```

### 2. 获取用户 (GetUser)

**请求:**
```bash
curl http://localhost:8088/v1/users/1
```

**成功响应:**
```json
{
  "error_code": 0,
  "error_msg": "success",
  "data": {
    "result": {
      "name": "users/1",
      "email": "alice@example.com",
      "display_name": "Alice Smith",
      "phone_number": "+1234567890",
      "create_time": "2025-12-17T10:00:00Z",
      "update_time": "2025-12-17T10:00:00Z",
      "is_active": true
    }
  }
}
```

**错误响应 (不存在):**
```json
{
  "error_code": 404,
  "error_msg": "user users/999 not found",
  "data": null
}
```

### 3. 列出用户 (ListUsers)

**请求:**
```bash
curl "http://localhost:8088/v1/users?page_size=10"
```

**成功响应:**
```json
{
  "error_code": 0,
  "error_msg": "success",
  "data": {
    "result": {
      "users": [
        {
          "name": "users/1",
          "email": "alice@example.com",
          "display_name": "Alice Smith",
          "is_active": true
        },
        {
          "name": "users/2",
          "email": "bob@example.com",
          "display_name": "Bob Johnson",
          "is_active": true
        }
      ],
      "next_page_token": "10",
      "total_size": 25
    }
  }
}
```

### 4. 更新用户 (UpdateUser)

**请求:**
```bash
curl -X PATCH http://localhost:8088/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "user": {
      "name": "users/1",
      "display_name": "Alice Johnson",
      "phone_number": "+9876543210"
    }
  }'
```

**成功响应:**
```json
{
  "error_code": 0,
  "error_msg": "success",
  "data": {
    "result": {
      "name": "users/1",
      "email": "alice@example.com",
      "display_name": "Alice Johnson",
      "phone_number": "+9876543210",
      "create_time": "2025-12-17T10:00:00Z",
      "update_time": "2025-12-17T10:15:00Z",
      "is_active": true
    }
  }
}
```

### 5. 删除用户 (DeleteUser)

**请求:**
```bash
curl -X DELETE http://localhost:8088/v1/users/1
```

**成功响应:**
```json
{
  "error_code": 0,
  "error_msg": "success",
  "data": null
}
```

### 6. 批量获取用户 (BatchGetUsers)

**请求:**
```bash
curl "http://localhost:8088/v1/users:batchGet?names=users/1&names=users/2&names=users/3"
```

**成功响应:**
```json
{
  "error_code": 0,
  "error_msg": "success",
  "data": {
    "result": {
      "users": [
        {
          "name": "users/1",
          "email": "alice@example.com",
          "display_name": "Alice Smith",
          "is_active": true
        },
        {
          "name": "users/2",
          "email": "bob@example.com",
          "display_name": "Bob Johnson",
          "is_active": true
        }
      ]
    }
  }
}
```

## gRPC 使用示例

使用 grpcurl 测试 gRPC 接口：

```bash
# 创建用户
grpcurl -plaintext -d '{
  "user": {
    "email": "alice@example.com",
    "display_name": "Alice Smith"
  }
}' localhost:9099 api.v1.UserService/CreateUser

# 获取用户
grpcurl -plaintext -d '{
  "name": "users/1"
}' localhost:9099 api.v1.UserService/GetUser
```

## 客户端错误处理

### JavaScript/TypeScript 示例

```typescript
interface CommonResponse<T = any> {
  error_code: number;
  error_msg: string;
  data: {
    result?: T;
  } | null;
}

async function createUser(userData: any): Promise<any> {
  const response = await fetch('http://localhost:8088/v1/users', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ user: userData }),
  });

  const result: CommonResponse = await response.json();

  if (result.error_code !== 0) {
    throw new Error(`API Error ${result.error_code}: ${result.error_msg}`);
  }

  return result.data?.result;
}
```

### Go 客户端示例

```go
type CommonResponse struct {
    ErrorCode int32                  `json:"error_code"`
    ErrorMsg  string                 `json:"error_msg"`
    Data      map[string]interface{} `json:"data"`
}

func createUser(email, displayName string) (*User, error) {
    reqBody := map[string]interface{}{
        "user": map[string]string{
            "email":        email,
            "display_name": displayName,
        },
    }

    jsonData, _ := json.Marshal(reqBody)
    resp, err := http.Post(
        "http://localhost:8088/v1/users",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result CommonResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    if result.ErrorCode != 0 {
        return nil, fmt.Errorf("API error %d: %s", result.ErrorCode, result.ErrorMsg)
    }

    // Extract user from data.result
    userData := result.Data["result"]
    // ... parse user data

    return user, nil
}
```

## 注意事项

1. **错误处理**: 始终检查 `error_code` 字段，不要仅依赖 HTTP 状态码
2. **数据解析**: 响应数据在 `data.result` 字段中，需要根据具体 API 解析相应的结构
3. **类型安全**: 使用 `google.protobuf.Struct` 提供了灵活性，但客户端需要自行确保类型安全
4. **幂等性**: 某些操作（如 DELETE）在资源不存在时仍返回成功（error_code=0）

## Swagger 文档

访问 http://localhost:8088/swagger/ 查看完整的交互式 API 文档。
