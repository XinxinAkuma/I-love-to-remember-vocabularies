package AI

import (
	"I-love-to-remember-vocabularies/model"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

func ConvertToText(ques *model.Mode, start, end int) string {
	var text string
	for _, item := range ques.List[start:end] {
		text += item.Title + "\nA. " + item.AnswerA + "\nB. " + item.AnswerB + "\nC. " + item.AnswerC + "\nD. " + item.AnswerD + "\n"
	}
	return text
}

func GetAnswers(ques *model.Mode) (*model.Result, error) {
	err := godotenv.Load("./secret.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apikey := os.Getenv("API_KEY")
	apiUrl := os.Getenv("API_URL")
	timeStr := os.Getenv("TIME")

	var allAnswers = make([]model.Answer, 100)
	var res model.Result
	var cnt int
	res.List = make([]model.Answer, 100)
	for i := 0; i < len(ques.List); i += 5 {
		end := i + 5
		if end > len(ques.List) {
			end = len(ques.List)
		}

		text := ConvertToText(ques, i, end)
		messages := []model.Message{
			{Role: "user", Content: "我会给你5题翻译题，需要你输出每个翻译题正确答案的选项，给我返回格式为`x-x-x-x-x`的答案，x表示正确答案的选项，不需要给出额外输出。题目如下：" + text},
		}

		reqBody := model.RequestBody{
			Model:    "gpt-4o",
			Messages: messages,
			Stream:   false,
		}

		body, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("error marshalling request body: %w", err)
		}

		req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(body))
		if err != nil {
			return nil, fmt.Errorf("error creating request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+apikey)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error executing request: %w", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("failed to get answer: status code %d", resp.StatusCode)
		}

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %w", err)
		}

		var response model.ResponseBody
		err = json.Unmarshal(respBody, &response)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling response body: %w", err)
		}

		fmt.Println("response", response.Choices[0].Message.Content)

		reStr, exist := Rematch(response.Choices[0].Message.Content)

		if exist {
			//split := strings.Split(reStr, "-")
			for k := 0; k < 5; k++ {
				if (k+1 > len(reStr) || len(reStr) == 0) && cnt <= len(ques.List) {
					allAnswers[cnt] = model.Answer{
						Input:         "A",
						PaperDetailId: ques.List[cnt].PaperDetailId,
					}
					cnt++
					continue
				}
				allAnswers[cnt] = model.Answer{
					Input:         reStr[k],
					PaperDetailId: ques.List[cnt].PaperDetailId,
				}
				fmt.Println("paperDetailedID: ", ques.List[cnt].PaperDetailId, "input: ", reStr[k])
				cnt++
			}
		} else {
			for k := 0; k < 5; k++ {
			}
		}
		time.Sleep(2 * time.Second)
	}
	res.PaperId = ques.PaperId
	res.Type = ques.Type
	res.List = allAnswers

	fmt.Println("-----------------------------------------------------------")
	fmt.Println("res: ", res)
	fmt.Println("please wait for" + timeStr + " seconds to submit the paper")
	subTime, _ := strconv.Atoi(timeStr)
	time.Sleep(time.Duration(subTime) * time.Second)
	fmt.Println("-----------------------------------------------------------")
	fmt.Println("submitting is ok ....")
	return &res, nil
}
func Rematch(resp string) ([]string, bool) {
	//pattern := `[A-D]-[A-D]-[A-D]-[A-D]-[A-D]`
	pattern := `[A-D]`

	re := regexp.MustCompile(pattern)

	matches := re.FindAllString(resp, -1)

	//for _, match := range matches {
	//	fmt.Println("match: ", match)
	//}
	if matches == nil {
		return nil, false
	}

	return matches, true
}
