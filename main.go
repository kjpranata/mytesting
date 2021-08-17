package main

import (
	//"context"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

type Config struct {
	ApiNodes []string `json:"apiNodes"`
	Sleep    int      `json:"sleep"`
	Bot      string   `json:"botApiKey"`
}

func configLoader(fileName string) (Config, error) {
	var config Config
	//Open File and Load it
	configFile, err := os.Open("config.json")
	if err != nil {
		return config, err
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	return config, err
}

var client []*sdk.Client
var conf []*sdk.Config

func init() {
	var err error
	config, _ := configLoader("config.json")
	for i := 0; i < len(config.ApiNodes); i++ {
		conf, err = sdk.NewConfig(context.Background(), []string{config.ApiNodes[i]})
		if err != nil {
			panic(err)
		}
		client[i] = sdk.NewClient(nil, conf[i])
	}
}

func main() {
	config, _ := configLoader("config.json")
	fmt.Println("Acturus: " + config.ApiNodes[0])
	fmt.Println("Aldebaran: " + config.ApiNodes[1])
	fmt.Println("Big Calvin: " + config.ApiNodes[2])
	fmt.Println("Bot: " + config.Bot)
	fmt.Println("Sleep: ", config.Sleep)
}
