package main

import (
	"context"
	"fmt"

	"example.com/main/application"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	checkError(err)

	app := application.New()

	err = app.Start(context.TODO())
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
