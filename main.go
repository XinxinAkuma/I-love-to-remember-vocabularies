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

	croypto, execution, session, encryptedPassword := work.GetNecessary("xxxxxx")

	tokens := work.Login("xxxxxx", croypto, execution, session, encryptedPassword)

	weekStr := os.Getenv("WEEK")

	week, _ := strconv.Atoi(weekStr)
	mode := os.Getenv("MODE")
	work.Final(tokens[1], week, mode)
}
