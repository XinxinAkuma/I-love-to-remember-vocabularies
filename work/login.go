package work

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func Login(studentId string, croypto string, execution string, session string, encryptedPassword string) []string {
	URL := "https://sso.hdu.edu.cn/login"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	var tokens []string

	client := http.Client{
		Transport: tr,
		// 禁用自动重定向，这样我们可以手动处理每一次重定向
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	formData := url.Values{}
	formData.Set("username", studentId)
	formData.Set("type", "UsernamePassword")
	formData.Set("_eventId", "submit")
	formData.Set("geolocation", "")
	formData.Set("execution", execution)
	formData.Set("captcha_code", "")
	formData.Set("croypto", croypto)
	formData.Set("password", encryptedPassword)

	req, err := http.NewRequest("POST", URL, strings.NewReader(formData.Encode()))
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return tokens
	}

	req.Header.Set("Host", "sso.hdu.edu.cn")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"134\", \"Not:A-Brand\";v=\"24\", \"Google Chrome\";v=\"134\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"Windows\"")
	req.Header.Set("Origin", "https://sso.hdu.edu.cn")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Referer", "https://sso.hdu.edu.cn/login?service=https:%2F%2Fskl.hdu.edu.cn%2Fapi%2Fcas%2Flogin%3Fstate%3DhF1wAuJeZqdVwgBziGX%26index%3D&state=hF1wAuJeZqdVwgBziGX")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cookie", session)

	currentReq := req
	redirectCount := 0
	maxRedirects := 10

	for redirectCount < maxRedirects {
		// 发送当前请求
		resp, err := client.Do(currentReq)
		if err != nil {
			fmt.Printf("请求 #%d 失败: %v\n", redirectCount+1, err)
			return tokens
		}

		tokenValue := resp.Header.Get("X-Auth-Token")
		if tokenValue != "" {
			fmt.Printf("响应 #%d 找到X-Auth-Token: %s\n", redirectCount+1, tokenValue)
			tokens = append(tokens, tokenValue)
		}

		// 如果不是重定向状态码，结束循环
		if resp.StatusCode < 300 || resp.StatusCode >= 400 {
			err := resp.Body.Close()
			if err != nil {
				return nil
			}
			break
		}

		// 获取重定向URL
		location := resp.Header.Get("Location")
		if location == "" {
			err := resp.Body.Close()
			if err != nil {
				return nil
			}
			fmt.Printf("响应 #%d 是重定向但没有Location头\n", redirectCount+1)
			return tokens
		}

		// 解析重定向URL
		redirectURL, err := url.Parse(location)
		if err != nil {
			err := resp.Body.Close()
			if err != nil {
				return nil
			}
			fmt.Printf("解析重定向URL失败: %v\n", err)
			return tokens
		}

		// 如果是相对URL，转换为绝对URL
		if !redirectURL.IsAbs() {
			redirectURL = currentReq.URL.ResolveReference(redirectURL)
		}

		// 创建新的请求用于重定向
		newReq, err := http.NewRequest("GET", redirectURL.String(), nil)
		if err != nil {
			err := resp.Body.Close()
			if err != nil {
				return nil
			}
			fmt.Printf("创建重定向请求失败: %v\n", err)
			return tokens
		}

		// 复制原始请求的头信息到新请求
		for key, values := range currentReq.Header {
			// 不复制一些特定的头
			if key != "Content-Type" && key != "Content-Length" && key != "Host" {
				for _, value := range values {
					newReq.Header.Add(key, value)
				}
			}
		}

		newReq.Host = redirectURL.Host

		for _, cookie := range resp.Cookies() {
			newReq.AddCookie(cookie)
		}

		err = resp.Body.Close()
		if err != nil {
			return nil
		}

		currentReq = newReq
		redirectCount++
	}

	return tokens
}
