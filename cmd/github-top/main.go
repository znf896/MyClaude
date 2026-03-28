package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/zhangzhanghaimin/myclaude/githubtrending"
)

func main() {
	// 命令行参数
	var (
		count     = flag.Int("count", 10, "Number of top projects to fetch")
		minStars  = flag.Int("min-stars", 1000, "Minimum stars required")
		language  = flag.String("language", "", "Filter by programming language (e.g. go, python)")
		query     = flag.String("query", "", "Search keywords (e.g. 'ai', 'kubernetes', 'machine learning')")
		sortBy    = flag.String("sort-by", githubtrending.SortByStars, "Sort by: stars/forks/updated")
		order     = flag.String("order", githubtrending.OrderDesc, "Sort order: desc/asc")
		outputDir = flag.String("output", "./output", "Output directory for markdown file")
		tokenEnv  = flag.String("token-env", "GITHUB_TOKEN", "Environment variable name for GitHub token")
		fetchReadme = flag.Bool("fetch-readme", true, "Fetch README content as article")
	)

	flag.Parse()

	// 获取GitHub Token
	token := os.Getenv(*tokenEnv)

	fmt.Printf("Fetching top GitHub projects...\n")
	fmt.Printf("Parameters: count=%d, min-stars=%d, language=%s, query=%q, sort-by=%s\n\n",
		*count, *minStars, *language, *query, *sortBy)

	// 创建客户端
	var client *githubtrending.Client
	if token != "" {
		fmt.Printf("✓ Using GitHub token from environment variable %s\n\n", *tokenEnv)
		client = githubtrending.NewClient(githubtrending.WithToken(token))
	} else {
		fmt.Printf("⚠ No GitHub token found in %s, using unauthenticated request (rate limit 60/hour)\n\n", *tokenEnv)
		client = githubtrending.NewClient()
	}

	ctx := context.Background()

	opt := &githubtrending.TopOptions{
		Count:    *count,
		MinStars: *minStars,
		Language: *language,
		Query:    *query,
		SortBy:   *sortBy,
		Order:    *order,
	}

	var (
		result *githubtrending.SearchResult
		err    error
	)

	if *fetchReadme {
		result, err = client.GetTopProjectsWithREADME(ctx, opt)
	} else {
		result, err = client.GetTopProjects(ctx, opt)
	}

	if err != nil {
		log.Fatalf("Failed to fetch projects: %v", err)
	}

	// 打印结果到控制台
	fmt.Printf("✅ Done! Found %d projects matching criteria\n\n", result.TotalCount)
	fmt.Printf("Returned top %d projects:\n\n", len(result.Items))

	for i, repo := range result.Items {
		fmt.Printf("%2d. %-40s 🌟%7d 🍴%6d\n", i+1, repo.FullName, repo.Stars, repo.Forks)
		fmt.Printf("     Description: %s\n", truncate(repo.Description, 60))
		fmt.Printf("     URL: %s\n\n", repo.HTMLURL)
	}

	// 导出到文件
	filePath, err := githubtrending.ExportToFile(result, *outputDir)
	if err != nil {
		log.Fatalf("Failed to export to file: %v", err)
	}

	fmt.Printf("📝 Results saved to: %s\n", filePath)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
