package work

import (
	"I-love-to-remember-vocabularies/AI"
	"I-love-to-remember-vocabularies/get"
	"I-love-to-remember-vocabularies/model"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func RequestQuestion(token string, week int, mode string) *model.Mode {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 禁用证书验证
		},
	}
	client := http.Client{
		Transport: tr,
	}
	startTime := time.Now().UnixMilli()
	urls := fmt.Sprintf("https://skl.hdu.edu.cn/api/paper/new?type=%s&week=%d&startTime=%d", mode, week, startTime)
	req, err := http.NewRequest("GET", urls, nil)
	if err != nil {
		panic(err)
	}
	req.Header = get.GetHeaders(token)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		panic("bad status code: " + resp.Status)
	}
	q := new(model.Mode)
	if err := json.Unmarshal(body, q); err != nil {
		panic(err)
	}

	return q
}

func Submit(res *model.Result, token string) error {
	time.Sleep(5 * time.Second)
	b, err := json.Marshal(res)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", "https://skl.hdu.edu.cn/api/paper/save", bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header = get.GetHeaders(token)
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		fmt.Println("你不能短时间内自测/考试")
	}
	fmt.Println("提交试卷成功")
	fmt.Println("result", resp)

	return nil
}

func Final(token string, week int, mode string) {
	q := RequestQuestion(token, week, mode)
	res, err := AI.GetAnswers(q)
	if err != nil {
		panic(err)
	}
	err = Submit(res, token)
	if err != nil {
		return
	}
}
