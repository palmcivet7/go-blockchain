package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	MINING_DIFFICULTY = 3
	MINING_SENDER = "THE BLOCKCHAIN"
	MINING_REWARD = 1.000000000000000000
)

type Block struct {
	timestamp		int64
	nonce			int
	previousHash	[32]byte
	transactions	[]*Transaction
}

func NewBlock(nonce int, previousHash [32]byte, transactions	[]*Transaction) *Block {
	b := new(Block)
	b.timestamp = time.Now().UnixNano()
	b.nonce = nonce
	b.previousHash = previousHash
	b.transactions = transactions
	return b
}

func (b *Block) Print() {
	fmt.Printf("Timestamp		%d\n", b.timestamp)
	fmt.Printf("Nonce			%d\n", b.nonce)
	fmt.Printf("Previous_Hash		%x\n", b.previousHash)
	for _, t := range b.transactions {
		 t.Print() 
	}
}

func (b *Block) Hash() [32]byte { 
	m, _ := json.Marshal(b)
	return sha256.Sum256([]byte(m))
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct{
		Timestamp		int64			`json:"timestamp"`
		Nonce			int				`json:"nonce"`
		PreviousHash	[32]byte		`json:"previous_hash"`
		Transactions	[]*Transaction	`json:"transactions"`
	}{
		Timestamp: b.timestamp,
		Nonce: b.nonce,
		PreviousHash: b.previousHash,
		Transactions: b.transactions,
	})
}

type Blockchain struct {
	transactionPool		[]*Transaction
	chain				[]*Block
	blockchainAddress	string
}

func NewBlockchain(blockchainAddress string) *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.blockchainAddress = blockchainAddress
	bc.CreateBlock(0, b.Hash())
	return bc
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	return b
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain) - 1]
}

func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Block %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 59))
}

type Transaction struct {
	senderAddress		string
	receiverAddress		string
	value 				float32
}

func (bc *Blockchain) AddTransaction(sender string, receiver string, value float32) {
	t := NewTransaction(sender, receiver, value)
	bc.transactionPool = append(bc.transactionPool, t)
}

func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.transactionPool{
		transactions = append(transactions, NewTransaction(t.senderAddress, t.receiverAddress, t.value))
	}
	return transactions
}

func (bc *Blockchain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) bool {
	 zeros := strings.Repeat("0", difficulty)
	 guessBlock := Block{0, nonce, previousHash, transactions}
	 guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	 return guessHashStr[:difficulty] == zeros
}

func (bc *Blockchain) ProofOfWork() int {
	transactions  := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().Hash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULTY){
		nonce += 1
	}
	return nonce
}

func (bc *Blockchain) Mining() bool {
	bc.AddTransaction(MINING_SENDER, bc.blockchainAddress, MINING_REWARD)
	nonce := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, previousHash)
	log.Println("action=mining, status=success")
	return true
}

func (bc *Blockchain) CalculateTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32 = 0.0
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.value
			if blockchainAddress == t.receiverAddress {
				totalAmount += value
			}
			if blockchainAddress == t.senderAddress {
				totalAmount -= value
			}
		}
	}
	return totalAmount
}

func NewTransaction(sender string, receiver string, value float32) *Transaction {
	return &Transaction{sender, receiver, value}
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf(" sender_address		%s\n", t.senderAddress)
	fmt.Printf(" receiver_address	%s\n", t.receiverAddress)
	fmt.Printf(" value			%.18f\n", t.value)
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct{
		Sender		string		`json:"sender_address"`
		Receiver	string		`json:"receiver_address"`
		Value 		float32		`json:"value"`
	}{
		Sender:		t.senderAddress,
		Receiver:	t.receiverAddress,
		Value:		t.value,
	})
}

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	minerAddress := "miner_address"

	bc := NewBlockchain(minerAddress)
	// bc.Print()

	bc.AddTransaction("Alice", "Bob", 1.000000000000000000)
	bc.Mining()

	bc.AddTransaction("palm", "civet", 4.500000000000000000)
	bc.AddTransaction("Satoshi", "Vitalik", 8.706050400000000000)
	bc.Mining()
	// previousHash = bc.LastBlock().Hash() // these 3 lines are completed with bc.Mining()
	// nonce = bc.ProofOfWork()
	// bc.CreateBlock(nonce, previousHash)
	bc.Print()

	fmt.Printf("miner_address %.18f\n", bc.CalculateTotalAmount("miner_address"))
	fmt.Printf("Alice %.18f\n", bc.CalculateTotalAmount("Alice"))
	fmt.Printf("Bob %.18f\n", bc.CalculateTotalAmount("Bob"))
} 