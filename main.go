package main

// 判断url状态码，并输出
import (
	"bufio"
	"time"
	"fmt"
	"net/http"
	"os"
	"strings"
	"wesite-title/goscraper"
)

// 获取url title
func url_title(url string) string {
	var title string
	s, err := goscraper.Scrape(url, 5)
	if err != nil {
		fmt.Println(err)
	}
	if len(s.Preview.Title) > 0 {
		title = s.Preview.Title
	} else {
		title = "unknown"
	}
	return title
}

//getStatusCode 获取url状态码
func getStatusCode(url string) int {
	resp, err := http.Get(url)
	if err != nil {
		return -1
	}
	defer resp.Body.Close()
	return resp.StatusCode
}

// 读文件
func readFile(filePath string) ([]string, error) {
	var lines []string
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func main() {
	// 读取文件
	lines, err := readFile("urls.txt")
	if err != nil {
		fmt.Println(err)
	}

	// titles 标题
	var titles string
	//为明确每个变量的类型，避免输出出错
	// writer.Write([]string{"title", "url", "statuscode"})
	titles = "title,url,statuscode\n"
	var stringBuilder strings.Builder
	stringBuilder.WriteString(titles)

	// 写入csv

	var url1 string
	for _, url1 = range lines {
		if strings.HasPrefix(url1, "http://") || strings.HasPrefix(url1, "https://") {
			// fmt.Println(url1)
			// statusCode := getStatusCode(url1)
			var statusCode int = getStatusCode(url1)
			// Code := string(rune(statusCode))
			// 判断url是否存在
			timeout := time.Duration(1 * time.Second)
			client := http.Client{
				Timeout: timeout,
			}
			_, err := client.Get(url1)
			if err != nil {
				print("访问出错了！", err.Error())
				continue
			}
			var title string = url_title(url1)
			fmt.Println(title, url1, statusCode)

			dataString := fmt.Sprintf("%s, %s, %d\n", title, url1, statusCode) // 因为类型不同，输出会出错，用格式化输出试了一下可以输出。
			stringBuilder.WriteString(dataString)

		}
	}

	var url2 string
	for _, url2 = range lines {
		if !strings.HasPrefix(url2, "http://") && !strings.HasPrefix(url2, "https://") {
			url2 = "http://" + url2
			// fmt.Println(url2)
			// statusCode := getStatusCode(url2)
			var statusCode int = getStatusCode(url2)
			// 判断url是否存在
			timeout := time.Duration(1 * time.Second)
			client := http.Client{
				Timeout: timeout,
			}
			_, err := client.Get(url2)
			if err != nil {
				print("访问出错了！", err.Error())
				continue
			}
			
			var title string = url_title(url2)
			fmt.Println(title, url2, statusCode)
			// 输出到csv
			// writer.Write([]string{title, url2, string(rune(statusCode))})
			dataString := fmt.Sprintf("%s, %s, %d\n", title, url2, statusCode)
			stringBuilder.WriteString(dataString)
		}
	}
	// 输出到csv
	filename := "./test.csv"
	file, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModeAppend|os.ModePerm)
	dataString := stringBuilder.String()
	file.WriteString(dataString)
	file.Close()
}
