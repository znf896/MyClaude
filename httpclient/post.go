package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client 通用HTTP客户端封装
type Client struct {
	client         *http.Client
	defaultHeaders map[string]string
	timeout        time.Duration
}

// Option 配置选项函数类型
type Option func(*Client)

// WithTimeout 设置超时
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// WithDefaultHeader 添加默认请求头
func WithDefaultHeader(key, value string) Option {
	return func(c *Client) {
		c.defaultHeaders[key] = value
	}
}

// WithTransport 设置自定义传输
func WithTransport(transport http.RoundTripper) Option {
	return func(c *Client) {
		c.client.Transport = transport
	}
}

// NewClient 创建新的HTTP客户端
func NewClient(opts ...Option) *Client {
	c := &Client{
		defaultHeaders: make(map[string]string),
		timeout:        30 * time.Second,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	// 更新超时
	if c.timeout != c.client.Timeout {
		c.client.Timeout = c.timeout
	}

	return c
}

// Request POST请求参数
type Request struct {
	URL     string            `json:"-"`
	Headers map[string]string `json:"-"`
	Body    interface{}       `json:"body"`
	Query   map[string]string `json:"query,omitempty"`
}

// Response HTTP响应封装
type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

// PostJSON 发送JSON格式的POST请求
func (c *Client) PostJSON(ctx context.Context, req *Request) (*Response, error) {
	return c.PostJSONWithResponse(ctx, req, nil)
}

// PostJSONWithResponse 发送POST请求并解析响应到指定结构体
func (c *Client) PostJSONWithResponse(ctx context.Context, req *Request, resp interface{}) (*Response, error) {
	// 序列化请求体
	var bodyBytes []byte
	var err error

	if req.Body != nil {
		bodyBytes, err = json.Marshal(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	// 构建带查询参数的URL
	fullURL, err := buildURLWithQuery(req.URL, req.Query)
	if err != nil {
		return nil, err
	}

	// 创建请求
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置默认请求头
	for k, v := range c.defaultHeaders {
		httpReq.Header.Set(k, v)
	}

	// 设置JSON内容类型
	if _, ok := req.Headers["Content-Type"]; !ok {
		httpReq.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	// 设置自定义请求头
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	// 发送请求
	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer httpResp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	result := &Response{
		StatusCode: httpResp.StatusCode,
		Headers:    httpResp.Header,
		Body:       respBody,
	}

	// 如果需要解析响应到结构体
	if resp != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, resp); err != nil {
			return result, fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return result, nil
}

// PostForm 发送application/x-www-form-urlencoded格式的POST请求
func (c *Client) PostForm(ctx context.Context, req *Request) (*Response, error) {
	return c.PostFormWithResponse(ctx, req, nil)
}

// PostFormWithResponse 发送表单请求并解析响应
func (c *Client) PostFormWithResponse(ctx context.Context, req *Request, resp interface{}) (*Response, error) {
	formData := url.Values{}

	if req.Body != nil {
		switch body := req.Body.(type) {
		case map[string]string:
			for k, v := range body {
				formData.Add(k, v)
			}
		case url.Values:
			formData = body
		default:
			return nil, errors.New("unsupported body type for form request, expect map[string]string or url.Values")
		}
	}

	// 构建带查询参数的URL
	fullURL, err := buildURLWithQuery(req.URL, req.Query)
	if err != nil {
		return nil, err
	}

	// 创建请求
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置默认请求头
	for k, v := range c.defaultHeaders {
		httpReq.Header.Set(k, v)
	}

	// 设置表单内容类型
	if _, ok := req.Headers["Content-Type"]; !ok {
		httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	}

	// 设置自定义请求头
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	// 发送请求
	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer httpResp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	result := &Response{
		StatusCode: httpResp.StatusCode,
		Headers:    httpResp.Header,
		Body:       respBody,
	}

	// 如果需要解析响应到结构体
	if resp != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, resp); err != nil {
			return result, fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return result, nil
}

// PostRaw 发送原始字节的POST请求
func (c *Client) PostRaw(ctx context.Context, url string, body []byte, headers map[string]string) (*Response, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置默认请求头
	for k, v := range c.defaultHeaders {
		httpReq.Header.Set(k, v)
	}

	// 设置自定义请求头
	for k, v := range headers {
		httpReq.Header.Set(k, v)
	}

	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &Response{
		StatusCode: httpResp.StatusCode,
		Headers:    httpResp.Header,
		Body:       respBody,
	}, nil
}

// IsSuccess 检查响应是否成功(2xx状态码)
func (r *Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// String 获取响应体字符串
func (r *Response) String() string {
	return string(r.Body)
}

// UnmarshalJSON 解析响应JSON到结构体
func (r *Response) UnmarshalJSON(v interface{}) error {
	return json.Unmarshal(r.Body, v)
}

// DefaultClient 默认全局客户端
var DefaultClient = NewClient()

// PostJSON 快捷方法 - 使用默认客户端发送JSON POST请求
func PostJSON(ctx context.Context, req *Request) (*Response, error) {
	return DefaultClient.PostJSON(ctx, req)
}

// PostJSONWithResult 快捷方法 - 使用默认客户端发送请求并解析响应
func PostJSONWithResult[T any](ctx context.Context, req *Request) (*T, *Response, error) {
	var result T
	resp, err := DefaultClient.PostJSONWithResponse(ctx, req, &result)
	if err != nil {
		return nil, resp, err
	}
	return &result, resp, nil
}

// PostForm 快捷方法 - 使用默认客户端发送表单POST请求
func PostForm(ctx context.Context, req *Request) (*Response, error) {
	return DefaultClient.PostForm(ctx, req)
}

// buildURLWithQuery 构建带查询参数的URL
func buildURLWithQuery(baseURL string, query map[string]string) (string, error) {
	if query == nil || len(query) == 0 {
		return baseURL, nil
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	q := parsedURL.Query()
	for k, v := range query {
		q.Add(k, v)
	}
	parsedURL.RawQuery = q.Encode()

	return parsedURL.String(), nil
}
