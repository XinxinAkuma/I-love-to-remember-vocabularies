package main

import (
	"I-love-to-remember-vocabularies/work"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func main() {
	err := godotenv.Load("./secret.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	tokens := os.Getenv("TOKENS")
	weekStr := os.Getenv("WEEK")
	week, _ := strconv.Atoi(weekStr)
	mode := os.Getenv("MODE")
	work.Final(tokens, week, mode)
}
