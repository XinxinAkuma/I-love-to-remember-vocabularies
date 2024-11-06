package model

import "time"

type Mode struct {
	PaperId     string      `json:"paperId"`
	StudentName interface{} `json:"studentName"`
	StudentId   string      `json:"studentId"`
	Type        string      `json:"type"`
	SchoolYear  string      `json:"schoolYear"`
	Semester    string      `json:"semester"`
	Week        int         `json:"week"`
	StartTime   time.Time   `json:"startTime"`
	EndTime     interface{} `json:"endTime"`
	SchoolCode  string      `json:"schoolCode"`
	Major       string      `json:"major"`
	TotalTime   int         `json:"totalTime"`
	Mark        int         `json:"mark"`
	List        []ListItem  `json:"list"`
	FinalResult interface{} `json:"finalResult"`
	Status      interface{} `json:"status"`
}

type ListItem struct {
	PaperDetailId string      `json:"paperDetailId"`
	Title         string      `json:"title"`
	AnswerA       string      `json:"answerA"`
	AnswerB       string      `json:"answerB"`
	AnswerC       string      `json:"answerC"`
	AnswerD       string      `json:"answerD"`
	QuestionId    interface{} `json:"questionId"`
	QuestionNum   interface{} `json:"questionNum"`
	Answer        interface{} `json:"answer"`
	Input         interface{} `json:"input"`
	Level         int         `json:"level"`
	Cet           int         `json:"cet"`
	Right         interface{} `json:"right"`
}

type Result struct {
	PaperId string   `json:"paperId"`
	Type    string   `json:"type"`
	List    []Answer `json:"list"`
}

type Answer struct {
	Input         interface{} `json:"input"`
	PaperDetailId string      `json:"paperDetailId"`
}

type Question struct {
	PaperDetailId string `json:"paperDetailId"`
	Title         string `json:"title"`
	AnswerA       string `json:"answerA"`
	AnswerB       string `json:"answerB"`
	AnswerC       string `json:"answerC"`
	AnswerD       string `json:"answerD"`
}

type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ResponseBody struct {
	Choices []ChoiceBody `json:"choices"`
	Usage   Usage        `json:"usage"`
}

type ChoiceBody struct {
	Message Content `json:"message"`
}

type Content struct {
	Content string `json:"content"`
}

type Usage struct {
	TotalTokens int `json:"total_tokens"`
}
