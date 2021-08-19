package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"log"
	"os"
	"time"
)

type Config struct {
	ApiNodes     []string    `json:"apiNodes"`
	ApiNodesName []string    `json:"apiNodesName"`
	Sleep        int         `json:"sleep"`
	Bot          string      `json:"botApiKey"`
	Id           int64       `json:"channelID"`
	Alert        string      `json:"alertMsg"`
	Sync         *sdk.Height `json:"syncValue"`
	Final        *sdk.Height `json:"finalHeight"`
}

func configLoader(fileName string) (Config, error) {
	var config Config
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
var clientTemp []*sdk.Client
var apiName []string
var conf *sdk.Config
var height sdk.Height
var finalHeight sdk.Height
var forkedHeight sdk.Height

func init() {
	config, _ := configLoader("config.json")
	var err error
	bcHeightTemp := []*sdk.Height{}
	apiNameTemp := []string{}
	i := 0
	j := 0

	for i < len(config.ApiNodes) {
		conf, err = sdk.NewConfig(context.Background(), []string{config.ApiNodes[i]})
		if err != nil {
			fmt.Println(config.ApiNodesName[i] + " is Offline.")
		} else {
			clientTemp = append(clientTemp, sdk.NewClient(nil, conf))
			bcHeight, err := clientTemp[j].Blockchain.GetBlockchainHeight(context.Background())
			if err != nil {
				panic(err)
			}
			bcHeightTemp = append(bcHeightTemp, &bcHeight)
			apiNameTemp = append(apiNameTemp, config.ApiNodesName[i])
			j++
		}
		i++
	}

	heightHigh := *bcHeightTemp[0]
	for i := 0; i < len(bcHeightTemp); i++ {
		if heightHigh < *bcHeightTemp[i] {
			heightHigh = *bcHeightTemp[i]
		}
	}

	for i = 0; i < len(bcHeightTemp); i++ {
		syncTest := heightHigh - *bcHeightTemp[i]
		if syncTest < *config.Sync {
			client = append(client, clientTemp[i])
			apiName = append(apiName, apiNameTemp[i])
		}
	}
}

func getHeight() {
	config, _ := configLoader("config.json")
	heights := []sdk.Height{}
	for i := 0; i < len(client); i++ {
		bcheight, err := client[i].Blockchain.GetBlockchainHeight(context.Background())
		if err != nil {
			panic(err)
		}
		heights = append(heights, bcheight)
	}

	height = heights[0]
	for i := 0; i < len(heights); i++ {
		if height > heights[i] {
			height = heights[i]
		}
	}
	finalHeight = height - *config.Final
}

func bot() {
	config, _ := configLoader("config.json")

	bot, err := tgbotapi.NewBotAPI(config.Bot)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	msg := tgbotapi.NewMessage(config.Id, "Current Height : "+height.String()+"\n"+config.Alert+forkedHeight.String())

	bot.Send(msg)
}

func main() {
	config, _ := configLoader("config.json")
	getHeight()

	for {
		blocks := []*sdk.BlockInfo{}
		hashes := []*sdk.Hash{}
		var fork bool

		for i := 0; i < len(client); i++ {
			block, err := client[i].Blockchain.GetBlockByHeight(context.Background(), finalHeight)
			if err != nil {
				panic(err)
			}
			blocks = append(blocks, block)
		}

		for i := 0; i < len(blocks); i++ {
			if i == 0 {
				fmt.Println("Current Block Height :", height)
				fmt.Println("Final Block Height :", finalHeight)
			}
			fmt.Println("Block "+apiName[i]+" Hash :", blocks[i].BlockHash)
			hashes = append(hashes, blocks[i].BlockHash)
		}

		for i := 0; i < len(hashes)-1; i++ {
			if hashes[i].String() == hashes[i+1].String() {
				fork = false
			} else {
				fork = true
				break
			}
		}

		if fork == true {
			Red := "\033[31m"
			Reset := "\033[0m"
			fmt.Println(string(Red), "Chain Forked! Sending Alarm Now!", string(Reset))
			forkedHeight = finalHeight
			bot()
			fmt.Println(string(Red), "Alert Sent.", string(Reset))
			// break

			//comment this and uncomment break for proper working
			//I comment the break because Betelgeuse always have different hash, so it would directly stop if not commented
			time.Sleep(time.Duration(config.Sleep) * time.Second)
		} else {
			time.Sleep(time.Duration(config.Sleep) * time.Second)
		}

		height++
		finalHeight++
	}
}
