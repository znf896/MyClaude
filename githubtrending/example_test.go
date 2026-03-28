package githubtrending

import (
	"context"
	"fmt"
	"log"
)

// ExampleClient_GetTopProjects 示例：获取Star最多的Top项目
func ExampleClient_GetTopProjects() {
	// 创建客户端（可选配置Token）
	// client := NewClient(WithToken("your-github-token-here"))
	client := NewClient() // 不配置Token也能使用，速率限制60请求/小时

	ctx := context.Background()

	// 配置查询选项
	opt := &TopOptions{
		Count:    10,  // 获取Top 10
		MinStars: 1000, // 至少1000星
		Language: "go", // Go语言项目
	}

	// 获取Top项目
	result, err := client.GetTopProjects(ctx, opt)
	if err != nil {
		log.Fatalf("Failed to get top projects: %v", err)
	}

	fmt.Printf("Total projects matching: %d\n", result.TotalCount)
	fmt.Printf("Returned top %d projects:\n", len(result.Items))

	for i, repo := range result.Items {
		fmt.Printf("%2d. %-30s Stars: %8d  %s\n", i+1, repo.FullName, repo.Stars, repo.Description)
	}
}

// ExampleClient_GetTopProjectsWithREADME 示例：获取Top项目并包含README文章内容
func ExampleClient_GetTopProjectsWithREADME() {
	client := NewClient()
	ctx := context.Background()

	opt := &TopOptions{
		Count:    5,
		MinStars: 10000,
		Language: "go",
	}

	result, err := client.GetTopProjectsWithREADME(ctx, opt)
	if err != nil {
		log.Fatalf("Failed to get top projects: %v", err)
	}

	fmt.Printf("Got %d projects\n", len(result.Items))

	for i, repo := range result.Items {
		fmt.Printf("\n%d. %s\nStars: %d\nDescription: %s\n", i+1, repo.FullName, repo.Stars, repo.Description)
		if repo.README != nil {
			fmt.Printf("README length: %d bytes\n", len(repo.README.Content))
			if len(repo.README.Content) > 200 {
				fmt.Printf("README preview: %.200s...\n", repo.README.Content)
			} else {
				fmt.Printf("README: %s\n", repo.README.Content)
			}
		}
	}
}

// ExampleGetTopProjects 示例：使用默认客户端快捷方法
func ExampleGetTopProjects() {
	ctx := context.Background()

	opt := &TopOptions{
		Count: 10,
	}

	result, err := GetTopProjects(ctx, opt)
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}

	fmt.Printf("Found %d projects\n", len(result.Items))
}
