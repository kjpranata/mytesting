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
var conf *sdk.Config
var height sdk.Height
var forkedHeight *sdk.Height

func init() {
	config, _ := configLoader("config.json")

	var err error
	//This line is my testing for checking sync height
	// bcHeightTemp := []*sdk.Height{}
	i := 0

	for i < len(config.ApiNodes) {
		conf, err = sdk.NewConfig(context.Background(), []string{config.ApiNodes[i]})
		if err != nil {
			fmt.Println(config.ApiNodesName[i] + " is Offline.")
		} else {
			client = append(client, sdk.NewClient(nil, conf))

			//This line is my testing for checking sync height
			// bcHeight, err := clientTemp[i].Blockchain.GetBlockchainHeight(context.Background())
			// if err != nil {
			// 	panic(err)
			// }
			// bcHeightTemp = append(bcHeightTemp, &bcHeight)
		}
		i++
	}

	//This line is my testing for checking sync height
	// for i = 0; i < len(clientTemp)-1; i++ {
	// 	syncTest := *bcHeightTemp[i] - *bcHeightTemp[i+1]
	// 	if syncTest > *config.Sync && syncTest > 0 {
	// 		conf, err = sdk.NewConfig(context.Background(), []string{config.ApiNodes[i]})
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		client = append(clientTemp, sdk.NewClient(nil, conf))
	// 	}

	// 	// fmt.Println("Test")
	// }
}

func getHeight() {
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

	msg := tgbotapi.NewMessage(config.Id, config.Alert+forkedHeight.String())

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
			block, err := client[i].Blockchain.GetBlockByHeight(context.Background(), height)
			if err != nil {
				panic(err)
			}
			blocks = append(blocks, block)
		}

		for i := 0; i < len(blocks); i++ {
			if i == 0 {
				fmt.Println("Block Height :", height)
			}
			fmt.Println("Block "+config.ApiNodesName[i]+" Hash :", blocks[i].BlockHash)
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
			forkedHeight = &height
			bot()
			// break

			//comment this and uncomment break for proper working
			//I comment the break because Betelgeuse always have different hash, so it would directly stop if not commented
			time.Sleep(time.Duration(config.Sleep) * time.Second)
		} else {
			time.Sleep(time.Duration(config.Sleep) * time.Second)
		}

		height++
	}
}
