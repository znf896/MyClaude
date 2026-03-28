# MyClaude - Go工具库集合

本项目包含多个实用的Go工具库：

- [httpclient](#httpclient---通用http-post请求封装) - 通用HTTP请求客户端，支持JSON、表单、原始数据
- [githubtrending](#githubtrending---获取github-star最多的top项目) - 获取GitHub上Star最多的热门项目及其README文章

---

## httpclient - 通用HTTP POST请求封装

这是一个用Go语言实现的通用HTTP POST请求封装库，支持JSON、表单和原始数据请求。

### 功能特性

- 支持 `application/json` 格式POST请求
- 支持 `application/x-www-form-urlencoded` 表单POST请求
- 支持原始字节数据POST请求
- **新增** 支持GET请求
- 支持默认请求头全局设置
- 支持自定义超时
- 支持URL查询参数
- 支持泛型直接解析响应
- 支持context上下文传递
- 代码简洁，无外部依赖

## 安装

```bash
go get ./httpclient
```

## 快速开始

### 1. 使用默认客户端发送JSON请求

```go
package main

import (
	"context"
	"fmt"
	"your-module-path/httpclient"
)

func main() {
	ctx := context.Background()

	// 定义请求体
	req := &httpclient.Request{
		URL:  "https://api.example.com/users",
		Body: map[string]interface{}{
			"name":  "John Doe",
			"email": "john@example.com",
		},
	}

	// 使用默认客户端发送请求
	resp, err := httpclient.PostJSON(ctx, req)
	if err != nil {
		panic(err)
	}

	if resp.IsSuccess() {
		fmt.Printf("Response: %s\n", resp.String())
	}
}
```

### 2. 创建自定义客户端

```go
import (
	"time"
	"your-module-path/httpclient"
)

// 创建自定义客户端，设置超时和默认请求头
client := httpclient.NewClient(
	httpclient.WithTimeout(10 * time.Second),
	httpclient.WithDefaultHeader("Authorization", "Bearer your-token"),
	httpclient.WithDefaultHeader("User-Agent", "MyApp/1.0"),
)
```

### 3. 使用泛型直接解析响应

```go
// 定义响应结构体
type UserResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    User   `json:"data"`
}

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// 直接返回解析后的结果
result, resp, err := httpclient.PostJSONWithResult[UserResponse](ctx, req)
if err != nil {
	panic(err)
}

fmt.Printf("User ID: %d, Name: %s\n", result.Data.ID, result.Data.Name)
```

### 4. 发送表单请求

```go
req := &httpclient.Request{
	URL:  "https://example.com/login",
	Body: map[string]string{
		"username": "admin",
		"password": "secret",
	},
}

resp, err := client.PostForm(ctx, req)
```

### 5. 带查询参数和自定义请求头

```go
req := &httpclient.Request{
	URL:     "https://api.example.com/search",
	Query:   map[string]string{"keyword": "go", "page": "1"},
	Headers: map[string]string{"X-Request-ID": "abc123"},
	Body: map[string]interface{}{
		"filter": "active",
	},
}
```

## API 说明

### 结构体

- `Client` - HTTP客户端主结构体
- `Request` - 请求参数结构体
  - `URL` - 请求地址
  - `Headers` - 自定义请求头
  - `Body` - 请求体
  - `Query` - URL查询参数
- `Response` - 响应封装
  - `StatusCode` - HTTP状态码
  - `Headers` - 响应头
  - `Body` - 响应字节数组

### 方法

- `NewClient(opts ...Option)` - 创建新客户端
- `Client.PostJSON(ctx, req)` - 发送JSON POST请求
- `Client.PostJSONWithResponse(ctx, req, resp)` - 发送请求并解析响应
- `Client.PostForm(ctx, req)` - 发送表单POST请求
- `Client.PostFormWithResponse(ctx, req, resp)` - 发送表单请求并解析响应
- `Client.PostRaw(ctx, url, body, headers)` - 发送原始字节请求
- `Response.IsSuccess()` - 检查是否是2xx成功状态
- `Response.String()` - 获取响应字符串
- `Response.UnmarshalJSON(v)` - 解析响应JSON

### 全局快捷方法

- `httpclient.PostJSON(ctx, req)` - 使用默认客户端发送JSON请求
- `httpclient.PostJSONWithResult[T](ctx, req)` - 使用泛型直接获取解析结果
- `httpclient.PostForm(ctx, req)` - 使用默认客户端发送表单请求

## 选项配置

- `WithTimeout(duration)` - 设置请求超时
- `WithDefaultHeader(key, value)` - 添加默认请求头
- `WithTransport(transport)` - 设置自定义RoundTripper

---

## githubtrending - 获取GitHub Star最多的Top项目

通过GitHub API获取Star最多的热门项目，同时可获取项目README文章内容。纯标准库实现，复用httpclient发送请求。

### 功能特性

- 获取按Star排序的热门项目
- 可筛选语言、最小Star数量
- 可选获取每个项目的README文章内容（自动base64解码）
- 支持GitHub Token认证（提升API速率限制）
- 选项模式配置，无外部依赖

### 快速开始

```go
package main

import (
	"context"
	"fmt"

	"github.com/zhangzhanghaimin/myclaude/githubtrending"
)

func main() {
	ctx := context.Background()

	// 配置查询选项
	opt := &githubtrending.TopOptions{
		Count:    10,  // 返回Top 10
		MinStars: 1000, // 最小Star数
		Language: "go",  // 只看Go语言项目
	}

	// 使用默认客户端获取Top项目
	result, err := githubtrending.GetTopProjects(ctx, opt)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Total projects: %d\n", result.TotalCount)
	for i, repo := range result.Items {
		fmt.Printf("%2d. %-30s Stars: %8d  %s\n", i+1, repo.FullName, repo.Stars, repo.Description)
	}
}
```

### 获取带README文章内容

```go
// 同时获取项目和README文章内容
result, err := githubtrending.GetTopProjectsWithREADME(ctx, opt)
if err != nil {
	panic(err)
}

for _, repo := range result.Items {
	fmt.Printf("Project: %s\n", repo.FullName)
	if repo.README != nil {
		fmt.Printf("README content length: %d bytes\n", len(repo.README.Content))
		fmt.Printf("README preview: %.200s...\n", repo.README.Content)
	}
}
```

### 配置GitHub Token

GitHub API对未认证请求限制60请求/小时，认证后提升到5000请求/小时：

```go
client := githubtrending.NewClient(
	githubtrending.WithToken("your-github-personal-access-token"),
	githubtrending.WithTimeout(30 * time.Second),
)

result, err := client.GetTopProjects(ctx, opt)
```

### API说明

- `NewClient(opts ...Option)` - 创建新客户端
- `Client.GetTopProjects(ctx, opt)` - 获取Top项目列表
- `Client.GetRepositoryREADME(ctx, owner, repo)` - 获取单个项目README
- `Client.GetTopProjectsWithREADME(ctx, opt)` - 获取Top项目并同时获取README
- `WithToken(token)` - 配置GitHub Token
- `WithTimeout(timeout)` - 配置超时

---

## 许可证

MIT
