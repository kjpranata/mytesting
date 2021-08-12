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
	forkHeight        = sdk.Height(3840000)
)

var ignoreAddress = map[string]string{
	"XDUYWYA5J7L4GBHOU34IXWVSBGEIWPB4ZHBVJKSI": "Patrick",
}

func main() {
	// Write Full Arcturus Log
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

	duration := time.Duration(50) * time.Millisecond

	//create files
	f, err1 := os.Create("Arcturus.txt")
	if err1 != nil {
		log.Fatal(err1)
	}
	defer f.Close()

	for height := cur_height; height > forkHeight; height-- {
		// Get TransactionInfo's by block height
		transactions, err2 := client.Blockchain.GetBlockTransactions(context.Background(), height)
		if err2 != nil {
			fmt.Printf("Blockchain.GetBlockTransactions returned error: %s", err2)
			return
		}

		fmt.Printf("Getting txn at %v height\n", height)
		if len(transactions) > 0 {
			// f.WriteString(height.String() + "\n")
			for _, transaction := range transactions {
				_, ignore := ignoreAddress[transaction.GetAbstractTransaction().Signer.Address.Address]
				// if transaction.GetAbstractTransaction().Signer.Address.Address != "XDUYWYA5J7L4GBHOU34IXWVSBGEIWPB4ZHBVJKSI" {
				if !ignore {
					fmt.Printf("Txn hash found: %v\n", transaction.GetAbstractTransaction().TransactionHash)
					f.WriteString("\t" + transaction.GetAbstractTransaction().TransactionHash.String() + "\n")
				}
			}
		}

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
			if transaction.GetAbstractTransaction().Signer.Address.Address != "XDUYWYA5J7L4GBHOU34IXWVSBGEIWPB4ZHBVJKSI" {
				f.WriteString("\n" + height.String())
				f.WriteString(transaction.String())
			}
		}

		height++
		time.Sleep(duration)
	}
}
