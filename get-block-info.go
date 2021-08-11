package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

const (
	// Sirius api rest server
	baseUrlArcturus   = "http://arcturus.xpxsirius.io:3000"
	baseUrlBetelgeuse = "http://betelgeuse.xpxsirius.io:3000"
)

func main() {
	//Write Full Arcturus Log
	arcturus()

	//Write Full Betelgeuse Log
	betelgeuse()
}

func arcturus() {
	conf, err := sdk.NewConfig(context.Background(), []string{baseUrlArcturus})
	if err != nil {
		fmt.Printf("NewConfig returned error: %s", err)
		return
	}

	// Use the default http client
	client := sdk.NewClient(nil, conf)

	//get current height
	cur_height, err := client.Blockchain.GetBlockchainHeight(context.Background())
	if err != nil {
		fmt.Printf("Blockchain.GetBlockhainHeight returned error: %s", err)
		return
	}

	// height of block in blockchain
	height := sdk.Height(3840000)

	duration := time.Duration(50) * time.Millisecond

	//create files
	f, err1 := os.Create("Arcturus.txt")
	if err1 != nil {
		log.Fatal(err1)
	}
	defer f.Close()

	for height < cur_height {

		// Get TransactionInfo's by block height
		transactions, err2 := client.Blockchain.GetBlockTransactions(context.Background(), height)
		if err2 != nil {
			fmt.Printf("Blockchain.GetBlockTransactions returned error: %s", err2)
			return
		}

		for _, transaction := range transactions {
			f.WriteString("\n" + height.String())
			f.WriteString(transaction.String())
		}

		height++
		time.Sleep(duration)
	}
}

func betelgeuse() {
	conf, err := sdk.NewConfig(context.Background(), []string{baseUrlBetelgeuse})
	if err != nil {
		fmt.Printf("NewConfig returned error: %s", err)
		return
	}

	// Use the default http client
	client := sdk.NewClient(nil, conf)

	//get current height
	cur_height, err := client.Blockchain.GetBlockchainHeight(context.Background())
	if err != nil {
		fmt.Printf("Blockchain.GetBlockhainHeight returned error: %s", err)
		return
	}

	// height of block in blockchain
	height := sdk.Height(3840000)

	duration := time.Duration(50) * time.Millisecond

	//create files
	f, err1 := os.Create("Betelgeuse.txt")
	if err1 != nil {
		log.Fatal(err1)
	}
	defer f.Close()

	for height < cur_height {

		// Get TransactionInfo's by block height
		transactions, err2 := client.Blockchain.GetBlockTransactions(context.Background(), height)
		if err2 != nil {
			fmt.Printf("Blockchain.GetBlockTransactions returned error: %s", err2)
			return
		}

		for _, transaction := range transactions {
			f.WriteString("\n" + height.String())
			f.WriteString(transaction.String())
		}

		height++
		time.Sleep(duration)
	}
}
