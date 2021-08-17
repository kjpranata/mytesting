// package main

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
// )

// const (
// 	// Sirius api rest server
// 	baseUrlA = "http://arcturus.xpxsirius.io:3000"
// 	baseUrlB = "http://betelgeuse.xpxsirius.io:3000"
// 	baseUrlC = "http://bigcalvin.xpxsirius.io:3000"
// )

// var clientA *sdk.Client
// var clientB *sdk.Client
// var clientC *sdk.Client
// var conf *sdk.Config

// func init() {
// 	var err error

// 	conf, err = sdk.NewConfig(context.Background(), []string{baseUrlA})
// 	if err != nil {
// 		panic(err)
// 	}
// 	clientA = sdk.NewClient(nil, conf)

// 	conf, err = sdk.NewConfig(context.Background(), []string{baseUrlB})
// 	if err != nil {
// 		panic(err)
// 	}
// 	clientB = sdk.NewClient(nil, conf)

// 	conf, err = sdk.NewConfig(context.Background(), []string{baseUrlC})
// 	if err != nil {
// 		panic(err)
// 	}
// 	clientC = sdk.NewClient(nil, conf)
// }

// func main() {
// 	curHeightA, err := clientA.Blockchain.GetBlockchainHeight(context.Background())
// 	if err != nil {
// 		panic(err)
// 	}

// 	curHeightB, err := clientB.Blockchain.GetBlockchainHeight(context.Background())
// 	if err != nil {
// 		panic(err)
// 	}

// 	curHeightC, err := clientC.Blockchain.GetBlockchainHeight(context.Background())
// 	if err != nil {
// 		panic(err)
// 	}

// 	var height sdk.Height

// 	heightArr := [3]sdk.Height{curHeightA, curHeightB, curHeightC}
// 	height = heightArr[0]
// 	for _, i := range heightArr {
// 		if i < height {
// 			height = i
// 		}
// 	}

// 	for {
// 		blockA, err := clientA.Blockchain.GetBlockByHeight(context.Background(), height)
// 		if err != nil {
// 			panic(err)
// 		}

// 		blockB, err := clientB.Blockchain.GetBlockByHeight(context.Background(), height)
// 		if err != nil {
// 			panic(err)
// 		}

// 		blockC, err := clientC.Blockchain.GetBlockByHeight(context.Background(), height)
// 		if err != nil {
// 			panic(err)
// 		}

// 		fmt.Printf("Block height: %v\n", height)
// 		fmt.Printf("\tblock A hash: %s\n", blockA.BlockHash)
// 		fmt.Printf("\tblock B hash: %s\n", blockB.BlockHash)
// 		fmt.Printf("\tblock C hash: %s\n\n", blockC.BlockHash)

// 		Red := "\033[31m"
// 		Reset := "\033[0m"

// 		if blockA.BlockHash != blockB.BlockHash || blockA.BlockHash != blockC.BlockHash {
// 			fmt.Println(string(Red), "Chain Forked !", string(Reset))
// 			break
// 		}
// 		height--

// 		time.Sleep(60 * time.Second)
// 	}
// }
