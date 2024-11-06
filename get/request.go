package get

import (
	"crypto/rand"
	"net/http"
)

func GetHeaders(token string) http.Header {
	ticket := generateTicket(21) // 自定义函数，模拟 JavaScript 中的 ticket 函数
	headers := http.Header{}
	headers.Set("Host", "skl.hdu.edu.cn")
	headers.Set("Connection", "keep-alive")
	headers.Set("Accept", "application/json, text/plain, */*")
	headers.Set("X-Auth-Token", token)
	headers.Set("User-Agent", "Mozilla/5.0 (Linux; U; Android 14; zh-CN; V2403A Build/UP1A.231005.007) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/100.0.4896.58 UWS/5.12.4.0 Mobile Safari/537.36 AliApp(DingTalk/7.6.20)")
	headers.Set("Skl-Ticket", ticket)
	headers.Set("Origin", "https://skl.hduhelp.com")
	headers.Set("X-Requested-With", "com.alibaba.android.rimet")
	headers.Set("Sec-Fetch-Site", "cross-site")
	headers.Set("Sec-Fetch-Mode", "cors")
	headers.Set("Sec-Fetch-Dest", "empty")
	headers.Set("Referer", "https://skl.hduhelp.com/")
	headers.Set("Accept-Encoding", "gzip, deflate, br")
	headers.Set("Accept-Language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
	return headers
}

func generateTicket(length int) string {
	const NL = "useandom-26T198340PX75pxJACKVERYMINDBUSHWOLF_GQZbfghjklqvwyzrict"
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	ticket := make([]byte, length)
	for i, b := range bytes {
		ticket[i] = NL[b&63] // 与 `& 63` 相同，确保索引在 0-63 范围内
	}
	return string(ticket)
}
