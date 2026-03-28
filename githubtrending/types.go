package githubtrending

// SortBy 排序方式常量
const (
	SortByStars     = "stars"     // 按Star数排序（默认）
	SortByForks     = "forks"     // 按Fork数排序
	SortByUpdated   = "updated"   // 按更新时间排序
)

// Order 排序顺序常量
const (
	OrderDesc = "desc" // 降序（默认）
	OrderAsc  = "asc"  // 升序
)

// TopOptions 查询选项
type TopOptions struct {
	Count     int    // 返回数量，默认30，最大100
	MinStars  int    // 最小star数，默认1
	Language  string // 按语言筛选，可选
	Query     string // 搜索关键词（支持自定义主题，如"ai", "kubernetes"等）
	SortBy    string // 排序方式：stars/forks/updated，默认stars
	Order     string // 排序顺序：desc/asc，默认desc
}

// SearchResult GitHub搜索结果
type SearchResult struct {
	TotalCount int           `json:"total_count"`
	Incomplete bool          `json:"incomplete_results"`
	Items      []*Repository `json:"items"`
}

// Repository GitHub仓库信息
type Repository struct {
	ID                int64       `json:"id"`
	Name              string      `json:"name"`
	FullName          string      `json:"full_name"`
	Private           bool        `json:"private"`
	Owner             Owner       `json:"owner"`
	Description       string      `json:"description"`
	HTMLURL           string      `json:"html_url"`
	CloneURL          string      `json:"clone_url"`
	SSHURL            string      `json:"ssh_url"`
	Stars             int         `json:"stargazers_count"`
	Watchers          int         `json:"watchers_count"`
	Forks             int         `json:"forks_count"`
	OpenIssues        int         `json:"open_issues_count"`
	Language          string      `json:"language"`
	CreatedAt         string      `json:"created_at"`
	UpdatedAt         string      `json:"updated_at"`
	PushedAt          string      `json:"pushed_at"`
	License           License     `json:"license"`
	README            *README     `json:"readme,omitempty"` // README文章内容
}

// Owner 仓库所有者信息
type Owner struct {
	Login     string `json:"login"`
	ID        int64  `json:"id"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
	Type      string `json:"type"`
}

// License 许可证信息
type License struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	SpdxID string `json:"spdx_id"`
	URL    string `json:"url"`
	NodeID string `json:"node_id"`
}

// README README文章内容
type README struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	SHA         string `json:"sha"`
	Size        int    `json:"size"`
	HTMLURL     string `json:"html_url"`
	DownloadURL string `json:"download_url"`
	Content     string `json:"content"` // 自动解码后的文本内容
}
