package main

import (
	"context"
	"fmt"

	"example.com/main/application"
)

func main() {
	app := application.New()

	err := app.Start(context.TODO())
	if err != nil {
		fmt.Println(err)
	}
}
