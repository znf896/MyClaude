package httpclient

import (
	"context"
	"fmt"
	"log"
	"time"
)

// ExampleClient_PostJSON 示例：发送JSON格式POST请求
func ExampleClient_PostJSON() {
	// 创建自定义客户端
	client := NewClient(
		WithTimeout(10*time.Second),
		WithDefaultHeader("Authorization", "Bearer your-token-here"),
		WithDefaultHeader("User-Agent", "MyApp/1.0"),
	)

	// 定义请求体
	type RequestBody struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Age      int    `json:"age"`
	}

	// 定义响应结构体
	type ResponseBody struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			UserID int64 `json:"user_id"`
		} `json:"data"`
	}

	ctx := context.Background()

	// 构造请求
	req := &Request{
		URL: "https://api.example.com/users/create",
		Body: RequestBody{
			Username: "johndoe",
			Email:    "john@example.com",
			Age:      25,
		},
		Query: map[string]string{
			"version": "v1",
		},
		Headers: map[string]string{
			"X-Request-ID": "req-123456",
		},
	}

	// 方式1：直接获取响应，手动解析
	resp, err := client.PostJSON(ctx, req)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}

	if !resp.IsSuccess() {
		log.Printf("Request failed with status code: %d", resp.StatusCode)
		log.Printf("Response: %s", resp.String())
		return
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Response length: %d bytes\n", len(resp.Body))

	// 解析响应到结构体
	var result ResponseBody
	if err := resp.UnmarshalJSON(&result); err != nil {
		log.Fatalf("Failed to parse response: %v", err)
	}

	fmt.Printf("Response: %+v\n", result)

	// 方式2：直接解析响应到结构体（泛型版本使用默认客户端）
	result2, resp, err := PostJSONWithResult[ResponseBody](ctx, req)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	fmt.Printf("Result with generic: %+v\n", result2)
}

// ExampleClient_PostForm 示例：发送表单POST请求
func ExampleClient_PostForm() {
	client := NewClient(WithTimeout(5*time.Second))

	ctx := context.Background()

	// 表单数据使用map[string]string
	req := &Request{
		URL:  "https://example.com/login",
		Body: map[string]string{"username": "admin", "password": "secret"},
	}

	resp, err := client.PostForm(ctx, req)
	if err != nil {
		log.Fatalf("Form request failed: %v", err)
	}

	fmt.Printf("Form login status: %d\n", resp.StatusCode)
}

// ExamplePostJSON 示例：使用默认客户端快捷方法
func ExamplePostJSON() {
	ctx := context.Background()

	req := &Request{
		URL:  "https://jsonplaceholder.typicode.com/posts",
		Body: map[string]interface{}{"title": "foo", "body": "bar", "userId": 1},
	}

	resp, err := PostJSON(ctx, req)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}

	fmt.Printf("Status: %d, Response: %s\n", resp.StatusCode, resp.String())
}
