package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"os"
)

type Config struct {
	ApiNodes struct {
		Acturus   string `json:"acturus"`
		Aldebaran string `json:"aldebaran"`
		Bigcalvin string `json:"bigcalvin"`
	} `json:"apiNodes"`
	Sleep int    `json:"sleep"`
	Bot   string `json:"botApiKey"`
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

	conf, err = sdk.NewConfig(context.Background(), []string{config.ApiNodes.Acturus})
	if err != nil {
		panic(err)
	}
	client[0] = sdk.NewClient(nil, conf[0])
}

func main() {
	config, _ := configLoader("config.json")
	fmt.Println("Acturus: " + config.ApiNodes.Acturus)
	fmt.Println("Aldebaran: " + config.ApiNodes.Aldebaran)
	fmt.Println("Big Calvin: " + config.ApiNodes.Bigcalvin)
	fmt.Println("Bot: " + config.Bot)
	fmt.Println("Sleep: ", config.Sleep)
}
