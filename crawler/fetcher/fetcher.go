package fetcher

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
)

var proxyUrl = "http://localhost:1087"

// Fetch 调用FetchRC获取[]byte类型的response.
func Fetch(pageUrl string, proxy bool) ([]byte, error) {
	rc, err := FetchRC(pageUrl, proxy)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return ioutil.ReadAll(rc)
}

// FetchRC 使用http包 获取ReadCloser类型.
func FetchRC(pageUrl string, proxy bool) (io.ReadCloser, error) {
	client := &http.Client{}
	// 是否使用代理
	if proxy {
		urlproxy, _ := url.Parse(proxyUrl)
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(urlproxy),
		}
	}

	// 创建request, 发送请求, 并拦截错误.
	req, err := http.NewRequest("GET", pageUrl, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		return resp.Body, nil
	}

	defer resp.Body.Close()
	return nil, fmt.Errorf("wrong status code: %d", resp.StatusCode)

}

// FetchToFile 调用FetchRC 将内容写入文件.
func FetchToFile(fileUrl, filePath string) error {
	rc, err := FetchRC(fileUrl, true)
	if err != nil {
		return err
	}
	defer rc.Close()

	// 检查文件所在路径是否存在
	err = os.MkdirAll(path.Dir(filePath), os.ModePerm)
	if err != nil {
		return err
	}

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入文件
	_, err = io.Copy(file, rc)
	return nil
}