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
	// Write Full Arcturus Log
	arcturus()

	//Write Full Betelgeuse Log
	betelgeuse()

	//Write Duplicate Log
	duplicate()
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
			if transaction.GetAbstractTransaction().Signer.Address.Address != "XDUYWYA5J7L4GBHOU34IXWVSBGEIWPB4ZHBVJKSI" {
				f.WriteString("\n" + height.String())
				f.WriteString(transaction.String())
			}
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
			if transaction.GetAbstractTransaction().Signer.Address.Address != "XDUYWYA5J7L4GBHOU34IXWVSBGEIWPB4ZHBVJKSI" {
				f.WriteString("\n" + height.String())
				f.WriteString(transaction.String())
			}
		}

		height++
		time.Sleep(duration)
	}
}

func duplicate() {
	confA, err := sdk.NewConfig(context.Background(), []string{baseUrlArcturus})
	if err != nil {
		fmt.Printf("NewConfig returned error: %s", err)
		return
	}

	confB, err2 := sdk.NewConfig(context.Background(), []string{baseUrlBetelgeuse})
	if err2 != nil {
		fmt.Printf("NewConfig returned error: %s", err2)
		return
	}

	// Use the default http client
	clientArcturus := sdk.NewClient(nil, confA)
	clientBetelgeuse := sdk.NewClient(nil, confB)

	// //get current height Arcturus
	// cur_height_a, err3 := clientArcturus.Blockchain.GetBlockchainHeight(context.Background())
	// if err3 != nil {
	// 	fmt.Printf("Blockchain.GetBlockhainHeight returned error: %s", err3)
	// 	return
	// }

	// fmt.Println(cur_height_a)

	// //get current height Betelgeus
	// cur_height_b, err4 := clientBetelgeuse.Blockchain.GetBlockchainHeight(context.Background())
	// if err4 != nil {
	// 	fmt.Printf("Blockchain.GetBlockhainHeight returned error: %s", err4)
	// 	return
	// }

	// fmt.Println(cur_height_b)

	// height of block in blockchain
	height := sdk.Height(3840000)

	duration := time.Duration(50) * time.Millisecond

	// create files
	f, err5 := os.Create("Duplicate.txt")
	if err5 != nil {
		log.Fatal(err5)
	}

	defer f.Close()

	for height < 3860000 {
		transactionsA, err6 := clientArcturus.Blockchain.GetBlockTransactions(context.Background(), height)
		if err6 != nil {
			fmt.Printf("Blockchain.GetBlockTransactions returned error: %s", err6)
			return
		}

		transactionsB, err7 := clientBetelgeuse.Blockchain.GetBlockTransactions(context.Background(), height)
		if err7 != nil {
			fmt.Printf("Blockchain.GetBlockTransactions returned error: %s", err7)
			return
		}

		for _, transactionA := range transactionsA {
			for _, transactionB := range transactionsB {
				if transactionA.GetAbstractTransaction().TransactionHash == transactionB.GetAbstractTransaction().TransactionHash {
					f.WriteString(transactionA.String())
					f.WriteString(transactionB.String())
				}
			}
		}

		height++
		time.Sleep(duration)
	}
}
