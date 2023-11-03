package main

import (
	"fmt"
	"log"

	"github.com/palmcivet7/go-blockchain/block"
	"github.com/palmcivet7/go-blockchain/wallet"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	// w := wallet.NewWallet()
	// fmt.Println(w.PrivateKeyStr())
	// fmt.Println(w.PublicKeyStr())
	// fmt.Println(w.BlockchainAddress())
	walletMiner := wallet.NewWallet()
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()

	t := wallet.NewTransaction(walletA.PrivateKey(), walletA.PublicKey(), walletA.BlockchainAddress(), walletB.BlockchainAddress(), 1.234567891012345678)
	// fmt.Printf("signature %s\n", t.GenerateSignature())

	bc := block.NewBlockchain(walletMiner.BlockchainAddress())
	isAdded := bc.AddTransaction(
		walletA.BlockchainAddress(),
		walletB.BlockchainAddress(),
		1.234567891012345678,
		walletA.PublicKey(),
		t.GenerateSignature(),
	)
	fmt.Println("Added?", isAdded)

	bc.Mining()
	bc.Print()

	fmt.Printf("A %.18f\n", bc.CalculateTotalAmount(walletA.BlockchainAddress()))
	fmt.Printf("B %.18f\n", bc.CalculateTotalAmount(walletB.BlockchainAddress()))
	fmt.Printf("Miner %.18f\n", bc.CalculateTotalAmount(walletMiner.BlockchainAddress()))
}   