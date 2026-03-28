package githubtrending

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Fatal("NewClient returned nil")
	}

	// 测试选项模式
	token := "test-token"
	clientWithToken := NewClient(WithToken(token))
	if clientWithToken.token != token {
		t.Errorf("Expected token %s, got %s", token, clientWithToken.token)
	}

	baseURL := "https://custom-api.github.com"
	clientWithBaseURL := NewClient(WithBaseURL(baseURL))
	if clientWithBaseURL.baseURL != baseURL {
		t.Errorf("Expected baseURL %s, got %s", baseURL, clientWithBaseURL.baseURL)
	}
}

func TestBuildSearchQuery(t *testing.T) {
	tests := []struct {
		name     string
		options  *TopOptions
		expected string
	}{
		{
			name: "empty options",
			options: &TopOptions{
				MinStars: 0,
			},
			expected: "stars:>1",
		},
		{
			name: "with min stars 1000",
			options: &TopOptions{
				MinStars: 1000,
			},
			expected: "stars:>1000",
		},
		{
			name: "with go language",
			options: &TopOptions{
				MinStars: 1000,
				Language: "go",
			},
			expected: "stars:>1000 language:go",
		},
		{
			name: "with extra query",
			options: &TopOptions{
				MinStars: 1000,
				Language: "go",
				Query:    "kubernetes",
			},
			expected: "stars:>1000 language:go kubernetes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.options.buildSearchQuery()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestGetTopProjects_Mock(t *testing.T) {
	// 创建测试服务器模拟GitHub API
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查请求方法
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		// 检查路径
		if r.URL.Path != "/search/repositories" {
			t.Errorf("Expected path /search/repositories, got %s", r.URL.Path)
		}

		// 读取测试数据
		data, err := os.ReadFile("./testdata/search_response.json")
		if err != nil {
			t.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}))
	defer ts.Close()

	// 创建客户端使用测试服务器
	client := NewClient(WithBaseURL(ts.URL + "/"))
	ctx := context.Background()
	opt := &TopOptions{
		Count:    10,
		MinStars: 1000,
	}

	result, err := client.GetTopProjects(ctx, opt)
	if err != nil {
		t.Fatalf("GetTopProjects failed: %v", err)
	}

	if result.TotalCount != 123456 {
		t.Errorf("Expected total_count 123456, got %d", result.TotalCount)
	}

	if len(result.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(result.Items))
	}

	repo := result.Items[0]
	if repo.FullName != "octocat/Hello-World" {
		t.Errorf("Expected octocat/Hello-World, got %s", repo.FullName)
	}

	if repo.Stars != 1600 {
		t.Errorf("Expected 1600 stars, got %d", repo.Stars)
	}
}

func TestIntegration_GetTopProjects(t *testing.T) {
	// 这个测试需要网络访问和可选GITHUB_TOKEN
	token := os.Getenv("GITHUB_TOKEN")
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := NewClient()
	if token != "" {
		client = NewClient(WithToken(token))
	}

	ctx := context.Background()
	opt := &TopOptions{
		Count:    5,
		MinStars: 10000,
		Language: "go",
	}

	result, err := client.GetTopProjects(ctx, opt)
	if err != nil {
		t.Fatalf("GetTopProjects failed: %v", err)
	}

	if result.TotalCount <= 0 {
		t.Fatal("Expected at least one project")
	}

	if len(result.Items) != 5 {
		t.Errorf("Expected 5 items, got %d", len(result.Items))
	}

	// 检查第一个项目是否有必要字段
	first := result.Items[0]
	if first.ID <= 0 {
		t.Error("Expected valid ID")
	}
	if first.FullName == "" {
		t.Error("Expected non-empty FullName")
	}
	if first.Stars <= 0 {
		t.Error("Expected stars > 0")
	}

	t.Logf("Got %d projects, total %d", len(result.Items), result.TotalCount)
	for i, item := range result.Items {
		t.Logf("%d: %s - %d stars", i+1, item.FullName, item.Stars)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	// 从文件读取测试数据
	data, err := os.ReadFile("./testdata/search_response.json")
	if err != nil {
		t.Fatal(err)
	}

	var result SearchResult
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result.TotalCount != 123456 {
		t.Errorf("Expected total_count 123456, got %d", result.TotalCount)
	}

	if len(result.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(result.Items))
	}
}
