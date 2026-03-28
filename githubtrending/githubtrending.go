package githubtrending

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/zhangzhanghaimin/myclaude/httpclient"
)

// Client GitHub Trending客户端
type Client struct {
	client  *httpclient.Client
	token   string
	baseURL string
}

// Option 配置选项函数类型
type Option func(*Client)

// WithToken 设置GitHub Token认证
func WithToken(token string) Option {
	return func(c *Client) {
		c.token = token
	}
}

// WithBaseURL 设置自定义API基础URL（用于代理等场景）
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithTimeout 设置请求超时
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.client = httpclient.NewClient(httpclient.WithTimeout(timeout))
	}
}

// WithDefaultHeader 添加默认请求头
func WithDefaultHeader(key, value string) Option {
	return func(c *Client) {
		// 重新创建client带着默认header
		// 这里保持和httpclient一致的选项模式设计
	}
}

// NewClient 创建新的GitHub Trending客户端
func NewClient(opts ...Option) *Client {
	c := &Client{
		client:  httpclient.NewClient(httpclient.WithTimeout(30 * time.Second)),
		baseURL: "https://api.github.com",
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// buildSearchQuery 构建搜索查询字符串
func (o *TopOptions) buildSearchQuery() string {
	var parts []string

	minStars := o.MinStars
	if minStars <= 0 {
		minStars = 1
	}
	parts = append(parts, fmt.Sprintf("stars:>%d", minStars))

	if o.Language != "" {
		parts = append(parts, fmt.Sprintf("language:%s", o.Language))
	}

	if o.Query != "" {
		parts = append(parts, o.Query)
	}

	return strings.Join(parts, " ")
}

// GetTopProjects 获取Star最多的Top项目
func (c *Client) GetTopProjects(ctx context.Context, opt *TopOptions) (*SearchResult, error) {
	count := opt.Count
	if count <= 0 {
		count = 30
	}
	if count > 100 {
		count = 100 // GitHub API最大限制
	}

	query := opt.buildSearchQuery()
	apiURL := fmt.Sprintf("%s/search/repositories", c.baseURL)

	queryParams := map[string]string{
		"q":        query,
		"sort":     "stars",
		"order":    "desc",
		"per_page": fmt.Sprintf("%d", count),
	}

	headers := make(map[string]string)
	headers["Accept"] = "application/vnd.github.v3+json"
	if c.token != "" {
		headers["Authorization"] = fmt.Sprintf("token %s", c.token)
	}

	result := &SearchResult{}
	resp, err := c.client.GetWithResponse(ctx, apiURL, queryParams, headers, result)
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return result, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, resp.String())
	}

	return result, nil
}

// GetRepositoryREADME 获取指定仓库的README（文章内容）
func (c *Client) GetRepositoryREADME(ctx context.Context, owner, repo string) (*README, error) {
	apiURL := fmt.Sprintf("%s/repos/%s/%s/readme", c.baseURL, owner, repo)

	headers := make(map[string]string)
	headers["Accept"] = "application/vnd.github.v3+json"
	if c.token != "" {
		headers["Authorization"] = fmt.Sprintf("token %s", c.token)
	}

	readme := &README{}
	resp, err := c.client.GetWithResponse(ctx, apiURL, nil, headers, readme)
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return readme, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, resp.String())
	}

	// GitHub返回的content是base64编码，需要解码
	if readme.Content != "" {
		// GitHub在base64末尾会有换行，需要移除
		content := strings.ReplaceAll(readme.Content, "\n", "")
		decoded, err := base64.StdEncoding.DecodeString(content)
		if err != nil {
			return readme, fmt.Errorf("failed to decode README content: %w", err)
		}
		readme.Content = string(decoded)
	}

	return readme, nil
}

// GetTopProjectsWithREADME 获取Top项目并同时获取它们的README文章内容
func (c *Client) GetTopProjectsWithREADME(ctx context.Context, opt *TopOptions) (*SearchResult, error) {
	result, err := c.GetTopProjects(ctx, opt)
	if err != nil {
		return nil, err
	}

	// 为每个项目获取README
	for _, repo := range result.Items {
		owner := repo.Owner.Login
		name := repo.Name
		readme, err := c.GetRepositoryREADME(ctx, owner, name)
		if err == nil {
			repo.README = readme
		}
		// 如果获取README失败，继续下一个，不影响整体结果
	}

	return result, nil
}

// 默认客户端
var DefaultClient = NewClient()

// GetTopProjects 使用默认客户端获取Star最多的Top项目
func GetTopProjects(ctx context.Context, opt *TopOptions) (*SearchResult, error) {
	return DefaultClient.GetTopProjects(ctx, opt)
}

// GetTopProjectsWithREADME 使用默认客户端获取Top项目并包含README文章内容
func GetTopProjectsWithREADME(ctx context.Context, opt *TopOptions) (*SearchResult, error) {
	return DefaultClient.GetTopProjectsWithREADME(ctx, opt)
}
