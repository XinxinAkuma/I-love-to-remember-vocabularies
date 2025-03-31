package work

import (
	"bytes"
	"crypto/des"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"strings"
)

func GetNecessary(password string) (string, string, string, string) {
	URL := "https://sso.hdu.edu.cn/login?service=https:%2F%2Fskl.hdu.edu.cn%2Fapi%2Fcas%2Flogin%3Fstate%3DhF1wAuJeZqdVwgBziGX%26index%3D&state=hF1wAuJeZqdVwgBziGX"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 禁用证书验证
		},
	}
	client := http.Client{
		Transport: tr,
	}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Host", "sso.hdu.edu.cn")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"134\", \"Not:A-Brand\";v=\"24\", \"Google Chrome\";v=\"134\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"Windows\"")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// 读取响应体内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应体出错:", err)
		panic(err)
	}
	htmlContent := string(body)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		fmt.Println("解析 HTML 出错:", err)
		panic(err)
	}
	croypto := doc.Find("#login-croypto").Text()
	fmt.Println("croypto:", croypto)
	execution := doc.Find("#login-page-flowkey").Text()
	fmt.Println("execution:", execution)
	headerValues := resp.Header["Set-Cookie"]
	session := Session(headerValues[0])
	fmt.Println("session:", session)

	encryptedPassword, err := desEncrypt(croypto, password)
	if err != nil {
		fmt.Println("加密出错:", err)
		panic(err)
	}

	fmt.Println("加密后的密码:", encryptedPassword)

	return croypto, execution, session, encryptedPassword

}

func Session(cookie string) string {
	startIndex := strings.Index(cookie, "SESSION=")
	if startIndex == -1 {
		return ""
	}
	endIndex := strings.Index(cookie[startIndex:], ";")
	if endIndex == -1 {
		return cookie[startIndex:]
	}
	return cookie[startIndex : startIndex+endIndex]
}
func desEncrypt(key, plaintext string) (string, error) {
	// 将密钥进行 Base64 解码
	decodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", err
	}

	// 创建 DES 块
	block, err := des.NewCipher(decodedKey)
	if err != nil {
		return "", err
	}

	// 标准PKCS7填充
	padding := des.BlockSize - len(plaintext)%des.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	paddedPlaintext := append([]byte(plaintext), padtext...)

	// 手动实现 ECB 模式加密
	ciphertext := make([]byte, len(paddedPlaintext))
	for i := 0; i < len(paddedPlaintext); i += des.BlockSize {
		block.Encrypt(ciphertext[i:i+des.BlockSize], paddedPlaintext[i:i+des.BlockSize])
	}

	// 将加密结果进行 Base64 编码（只返回一次结果）
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
