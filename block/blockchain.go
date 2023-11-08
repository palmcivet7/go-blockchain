package block

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/palmcivet7/go-blockchain/utils"
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
		PreviousHash	string		`json:"previous_hash"`
		Transactions	[]*Transaction	`json:"transactions"`
	}{
		Timestamp: b.timestamp,
		Nonce: b.nonce,
		PreviousHash: fmt.Sprintf("%x", b.previousHash),
		Transactions: b.transactions,
	})
}

type Blockchain struct {
	transactionPool		[]*Transaction
	chain				[]*Block
	blockchainAddress	string
	port				uint16
}

func NewBlockchain(blockchainAddress string, port uint16) *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.blockchainAddress = blockchainAddress
	bc.CreateBlock(0, b.Hash())
	bc.port = port
	return bc
}

func (bc *Blockchain) TransactionPool() []*Transaction {
	return bc.transactionPool
}

func (bc *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct{
		Blocks []*Block	`json:"chains"`
	}{
		Blocks: bc.chain,
	})
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
	value 				float64
}

func (bc *Blockchain) CreateTransaction(
	sender string, receiver string, value float64, senderPublicKey *ecdsa.PublicKey, s *utils.Signature,
) bool {
	isTransacted := bc.AddTransaction(sender, receiver, value, senderPublicKey, s)

	// TODO
	// SYNC

	return isTransacted
}

func (bc *Blockchain) AddTransaction(
	sender string, receiver string, value float64, senderPublicKey *ecdsa.PublicKey, s *utils.Signature,
) bool {
	t := NewTransaction(sender, receiver, value)

	if sender == MINING_SENDER {
		bc.transactionPool = append(bc.transactionPool,  t)
		return true
	}

	if bc.VerifyTransactionSignature(senderPublicKey, s, t){
		// if bc.CalculateTotalAmount(sender) < value {
		// 	log.Println("ERROR: Not enough balance in wallet")
		// 	return false
		// }
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	} else {
		log.Println("ERROR: Verify Transaction")
	}
	return false
}

func (bc *Blockchain) VerifyTransactionSignature(
	senderPublicKey *ecdsa.PublicKey, s *utils.Signature, t *Transaction,
) bool {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	return ecdsa.Verify(senderPublicKey, h[:], s.R, s.S)
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
	bc.AddTransaction(MINING_SENDER, bc.blockchainAddress, MINING_REWARD, nil, nil)
	nonce := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, previousHash)
	log.Println("action=mining, status=success")
	return true
}

func (bc *Blockchain) CalculateTotalAmount(blockchainAddress string) float64 {
	var totalAmount float64 = 0.0
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

func NewTransaction(sender string, receiver string, value float64) *Transaction {
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
		Value 		float64		`json:"value"`
	}{
		Sender:		t.senderAddress,
		Receiver:	t.receiverAddress,
		Value:		t.value,
	})
}

type TransactionRequest struct {
	SenderAddress	*string 	`json:"sender_address"`
	ReceiverAddress *string 	`json:"receiver_address"`
	SenderPublicKey	*string		`json:"sender_public_key"`
	Value			*float64	`json:"value"`
	Signature 		*string		`json:"signature`

}

func (tr *TransactionRequest) Validate () bool {
	if tr.SenderAddress == nil ||
		tr.ReceiverAddress == nil ||
		tr.SenderPublicKey == nil ||
		tr.Value == nil ||
		tr.Signature == nil {
		return false
	}
	return true
}