package main

import "fmt"

func main() {
	fmt.Println("Dev Configuration")

	configuration := GetConfig()
	fmt.Println(configuration.DB_USERNAME)

	fmt.Println("Prod Configuration")

	configuration = GetConfig("prod")
	fmt.Println(configuration.DB_USERNAME)
}
