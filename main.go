package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

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
var conf *sdk.Config
var height sdk.Height

func init() {
	config, _ := configLoader("config.json")

	var err error
	for i := 0; i < len(config.ApiNodes); i++ {
		conf, err = sdk.NewConfig(context.Background(), []string{config.ApiNodes[i]})
		if err != nil {
			panic(err)
		}
		client = append(client, sdk.NewClient(nil, conf))
	}
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
			switch i {
			case 0:
				fmt.Println("Block Height:", height)
				fmt.Println("Block Aldebaran Hash  :", blocks[i].BlockHash)
			case 1:
				fmt.Println("Block Arcturus Hash   :", blocks[i].BlockHash)
			case 2:
				fmt.Println("Block Betelgeuse Hash :", blocks[i].BlockHash)
			case 3:
				fmt.Println("Block BigCalvin Hash  :", blocks[i].BlockHash)
			default:
			}

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
			break
		} else {
			time.Sleep(time.Duration(config.Sleep) * time.Second)
		}

		height++
	}
}
