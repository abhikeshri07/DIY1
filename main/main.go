package main

import (
	"fmt"
	"os"
)

func main() {
	a := App{}

	a.Initialize(os.Getenv("APP_DB,HOST"),
		os.Getenv("APP_DB_PORT"),
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"))
	fmt.Println("Database Connected and starting server at 8010")

	a.Run(os.Getenv("APP_DB_HOST"), "8010")
}
