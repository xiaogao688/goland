package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const (
	owner = "emn178" // 仓库所有者
	repo  = "js-md5" // 仓库名称
)

type GitRef struct {
	Ref    string `json:"ref"`
	NodeID string `json:"node_id"`
	URL    string `json:"url"`
	Object struct {
		SHA  string `json:"sha"`
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"object"`
}

func main() {
	// 获取分支引用
	branches, err := fetchRefs(fmt.Sprintf("https://api.github.com/repos/%s/%s/git/refs/heads", owner, repo))
	if err != nil {
		fmt.Println("Error fetching branches:", err)
		os.Exit(1)
	}

	// 获取标签引用
	tags, err := fetchRefs(fmt.Sprintf("https://api.github.com/repos/%s/%s/git/refs/tags", owner, repo))
	if err != nil {
		fmt.Println("Error fetching tags:", err)
		os.Exit(1)
	}

	// 打印分支引用
	fmt.Println("Branches:")
	for _, branch := range branches {
		fmt.Printf("%s: %s\n", branch.Ref, branch.Object.SHA)
	}

	// 打印标签引用
	fmt.Println("\nTags:")
	for _, tag := range tags {
		fmt.Printf("%s: %s\n", tag.Ref, tag.Object.SHA)
	}
}

func fetchRefs(url string) ([]GitRef, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 设置请求头，必要时添加认证信息
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	// 如果需要访问私人仓库，请取消注释以下行，并使用你的 GitHub 个人访问令牌
	// req.Header.Set("Authorization", "token YOUR_GITHUB_PERSONAL_ACCESS_TOKEN")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch refs: %s", resp.Status)
	}

	var refs []GitRef
	if err := json.NewDecoder(resp.Body).Decode(&refs); err != nil {
		return nil, err
	}

	return refs, nil
}
