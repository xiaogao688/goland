package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
	t2()
}

const url = "https://api.github.com/repos/emn178/js-md5/tags"

type Tag struct {
	Name       string `json:"name"`
	ZipballUrl string `json:"zipball_url"`
	TarballUrl string `json:"tarball_url"`
}

// 获取仓库所有的tag以及 tar地址和zip地址和commit
func t1() {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	var tags []Tag
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, tag := range tags {
		fmt.Printf("Name: %s\nZipball URL: %s\nTarball URL: %s\n\n", tag.Name, tag.ZipballUrl, tag.TarballUrl)
	}
}

// 向指定的 Git 仓库发送 ls-refs 命令
func t2() {
	// URL of the git repository
	url := "https://github.com/emn178/js-md5.git/git-upload-pack"

	// Payload for the ls-refs command
	payload := "0014command=ls-refs\n" +
		"0014agent=git/2.43.0\n" +
		"0016object-format=sha1\n" +
		"00010009peel\n" +
		"000c symrefs\n" +
		"000b unborn\n" +
		"001a ref-prefix refs/tags/\n" +
		"001b ref-prefix refs/heads/\n" +
		"0000"

	// Buffer to store the request body
	var requestBody bytes.Buffer
	requestBody.WriteString(payload)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(1)
	}

	// Set the required headers
	req.Header.Set("Content-Type", "application/x-git-upload-pack-request")
	req.Header.Set("Accept", "application/x-git-upload-pack-result")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Check if the response status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: received non-200 response code:", resp.Status)
		os.Exit(1)
	}

	// Parse the response
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "0000" {
			break
		}
		fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading response:", err)
		os.Exit(1)
	}

}
