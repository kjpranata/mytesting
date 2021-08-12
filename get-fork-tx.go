package main

import (
	"context"
	"fmt"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

const (
	// Sirius api rest server
	baseUrlA = "http://arcturus.xpxsirius.io:3000"
	baseUrlB = "http://betelgeuse.xpxsirius.io:3000"
)

var clientA *sdk.Client
var clientB *sdk.Client
var ctx context.Context
var conf *sdk.Config
var ignoreAddress = map[string]string{
	"XDUYWYA5J7L4GBHOU34IXWVSBGEIWPB4ZHBVJKSI": "Patrick",
}
var dups = map[string]bool{}

func init() {
	ctx = context.Background()
	var err error
	conf, err = sdk.NewConfig(ctx, []string{baseUrlA})
	if err != nil {
		panic(err)
	}

	clientA = sdk.NewClient(nil, conf)

	conf, err = sdk.NewConfig(ctx, []string{baseUrlB})
	if err != nil {
		panic(err)
	}

	clientB = sdk.NewClient(nil, conf)

}

func main() {
	curHeightA, err := clientA.Blockchain.GetBlockchainHeight(ctx)
	if err != nil {
		panic(err)
	}

	curHeightB, err := clientB.Blockchain.GetBlockchainHeight(ctx)
	if err != nil {
		panic(err)
	}

	var height sdk.Height

	if curHeightA > curHeightB {
		height = curHeightB
	} else {
		height = curHeightA
	}

	for {

		blockA, err := clientA.Blockchain.GetBlockByHeight(ctx, height)
		if err != nil {
			panic(err)
		}

		blockB, err := clientB.Blockchain.GetBlockByHeight(ctx, height)
		if err != nil {
			panic(err)
		}
		// fmt.Printf("%v ", height)
		// fmt.Printf("Block height: %v\n", height)
		// fmt.Printf("\tblock A hash: %s\n", blockA.BlockHash)
		// fmt.Printf("\tblock B hash: %s\n\n", blockB.BlockHash)

		if blockA.BlockHash.String() == blockB.BlockHash.String() {
			fmt.Printf("Block hash same at block height %v\n", height)
			break
		}

		if height <= curHeightA {
			checkBlkTx(clientA, height)
		}

		if height < curHeightB {
			checkBlkTx(clientB, height)
		}

		height--
	}

	fmt.Println(dups)
}

func checkBlkTx(client *sdk.Client, height sdk.Height) {
	blkTx, err := client.Blockchain.GetBlockTransactions(ctx, height)
	if err != nil {
		panic(err)
	}

	if len(blkTx) > 0 {
		for _, tx := range blkTx {
			_, ignore := ignoreAddress[tx.GetAbstractTransaction().Signer.Address.Address]
			if !ignore {
				hash := tx.GetAbstractTransaction().TransactionHash.String()
				if txExist(clientA, hash) && txExist(clientB, hash) {
					// fmt.Println(hash + " is duplicate txn")
					if dups[hash] != true {
						dups[hash] = true
					}
				} else {
					dups[hash] = false
				}
			}
		}
	}
}

func txExist(client *sdk.Client, hash string) bool {
	_, err := client.Transaction.GetTransaction(ctx, hash)
	if err != nil {
		return false
	}
	return true
}
