package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/zhangzhanghaimin/myclaude/githubtrending"
)

// APIResponse 统一API响应格式
type APIResponse struct {
	Code    int                     `json:"code"`
	Message string                  `json:"message"`
	Data    *githubtrending.SearchResult `json:"data,omitempty"`
}

// handleTop 处理获取top项目请求
func handleTop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析查询参数
	query := r.FormValue("query")
	language := r.FormValue("language")

	countStr := r.FormValue("count")
	count := 10
	if countStr != "" {
		if n, err := strconv.Atoi(countStr); err == nil && n > 0 && n <= 100 {
			count = n
		}
	}

	minStarsStr := r.FormValue("min_stars")
	minStars := 1000
	if minStarsStr != "" {
		if n, err := strconv.Atoi(minStarsStr); err == nil && n >= 0 {
			minStars = n
		}
	}

	sortBy := r.FormValue("sort_by")
	if sortBy == "" {
		sortBy = githubtrending.SortByStars
	}

	order := r.FormValue("order")
	if order == "" {
		order = githubtrending.OrderDesc
	}

	fetchReadmeStr := r.FormValue("fetch_readme")
	fetchReadme := true
	if fetchReadmeStr == "false" || fetchReadmeStr == "0" {
		fetchReadme = false
	}

	// 构建选项
	opt := &githubtrending.TopOptions{
		Count:    count,
		MinStars: minStars,
		Language: language,
		Query:    query,
		SortBy:   sortBy,
		Order:    order,
	}

	// 获取token从环境变量或header
	token := r.Header.Get("X-Github-Token")
	if token == "" {
		token = githubtrending.DefaultClientGetToken()
	}

	var client *githubtrending.Client
	if token != "" {
		client = githubtrending.NewClient(githubtrending.WithToken(token))
	} else {
		client = githubtrending.DefaultClient
	}

	ctx := context.Background()

	var result *githubtrending.SearchResult
	var err error

	if fetchReadme {
		result, err = client.GetTopProjectsWithREADME(ctx, opt)
	} else {
		result, err = client.GetTopProjects(ctx, opt)
	}

	if err != nil {
		response := APIResponse{
			Code:    500,
			Message: err.Error(),
			Data:    nil,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := APIResponse{
		Code:    200,
		Message: "success",
		Data:    result,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"service": "github-trending-api",
	})
}

func main() {
	addr := ":8080"
	log.Printf("Starting GitHub Trending API server on %s...\n", addr)
	log.Println("API endpoints:")
	log.Println("  GET  /api/health - health check")
	log.Println("  GET/POST /api/top - Get top GitHub projects")
	log.Println("    Query params:")
	log.Println("      query     - search keywords (e.g. 'ai artificial-intelligence')")
	log.Println("      language  - filter by programming language")
	log.Println("      count     - number of projects (default 10, max 100)")
	log.Println("      min_stars - minimum stars required (default 1000)")
	log.Println("      sort_by   - sort by: stars/forks/updated (default stars)")
	log.Println("      order     - sort order: desc/asc (default desc)")
	log.Println("      fetch_readme - fetch README content: true/false (default true)")
	log.Println("    Header:")
	log.Println("      X-Github-Token - optional GitHub token for higher rate limit")

	http.HandleFunc("/api/health", handleHealth)
	http.HandleFunc("/api/top", handleTop)

	log.Fatal(http.ListenAndServe(addr, nil))
}
