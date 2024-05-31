package main

import (
	"fmt"
	"regexp"
)

func main() {
	t2()
}

func t1() {
	// 示例文本，其中包含要替换的URL
	text := `Here are some URLs: http://example.com, https://example.org, http://example.com/page, and some other text.`

	// 要替换的URL
	urlToReplace := `http://example.com`

	// 定义正则表达式，匹配指定的 URL
	re := regexp.MustCompile(regexp.QuoteMeta(urlToReplace))

	// 使用 ReplaceAllString 方法进行替换
	replacedText := re.ReplaceAllString(text, "art://")

	fmt.Println("Before replacement:", text)
	fmt.Println("After replacement:", replacedText)
}

func t2() {
	// 示例文本，其中包含要处理的URL
	text := `Here are some URLs: http://example.com.git, https://example.org.git, and some other text ending with .git.`

	// 定义正则表达式，仅匹配字符串结尾的 .git
	re := regexp.MustCompile(`.git$`)

	// 使用 ReplaceAllString 方法进行替换，将字符串结尾的 .git 去掉
	replacedText := re.ReplaceAllString(text, "")

	fmt.Println("Before replacement:", text)
	fmt.Println("After replacement:", replacedText)
}
