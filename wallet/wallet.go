package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"github.com/palmcivet7/go-blockchain/utils"
)

type Wallet struct {
	privateKey        *ecdsa.PrivateKey
	publicKey         *ecdsa.PublicKey
	blockchainAddress string
}

func NewWallet() *Wallet {
	w := new(Wallet)
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	w.privateKey = privateKey
	w.publicKey = &w.privateKey.PublicKey

	h2 := sha256.New()
	h2.Write(w.publicKey.X.Bytes())
	h2.Write(w.publicKey.Y.Bytes())
	digest2 := h2.Sum(nil)

	h3 := sha256.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)

	vd4 := make([]byte, 33) // Adjusted size from 21 to 33
	vd4[0] = 0x00
	copy(vd4[1:], digest3[:])

	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)

	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)

	chsum := digest6[:4]

	dc8 := make([]byte, 37) // Adjusted size from 25 to 37
	copy(dc8[:33], vd4)     // Adjusted size from 21 to 33
	copy(dc8[33:], chsum)

	address := base58.Encode(dc8)
	w.blockchainAddress = address
	return w
}

func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

func (w *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", w.privateKey.D.Bytes())
}

func (w *Wallet) PublicKey() *ecdsa.PublicKey {
	return w.publicKey
}

func (w *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%064x%064x", w.publicKey.X.Bytes(), w.publicKey.Y.Bytes())
}

func (w *Wallet) BlockchainAddress() string {
	return w.blockchainAddress
}

func (w *Wallet) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct{
		PrivateKey			string	`json:"private_key"`
		PublicKey			string	`json:"public_key"`
		BlockchainAddress	string	`json:"blockchain_address"`
	}{
		PrivateKey: w.PrivateKeyStr(),
		PublicKey: w.PublicKeyStr(),
		BlockchainAddress: w.BlockchainAddress(),
	})
}

type Transaction struct {
	senderPrivateKey			*ecdsa.PrivateKey
	senderPublicKey				*ecdsa.PublicKey
	senderBlockchainAddress		string
	receiverBlockchainAddress	string
	value						float64
}

func NewTransaction(
	privateKey *ecdsa.PrivateKey,
	publicKey *ecdsa.PublicKey,
	sender string,
	receiver string,
	value float64,
) *Transaction {
	return &Transaction{
		privateKey, publicKey, sender, receiver, value,
	}
} 

func (t *Transaction) GenerateSignature() *utils.Signature {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	r, s, _ := ecdsa.Sign(rand.Reader, t.senderPrivateKey, h[:])
	return &utils.Signature{r, s}
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct{
		Sender		string	`json:"sender_address"`
		Receiver	string	`json:"receiver_address"`
		Value		float64	`json:"value"`
	}{
		Sender: t.senderBlockchainAddress,
		Receiver: t.receiverBlockchainAddress,
		Value: t.value,
	})
}

type TransactionRequest struct {
	SenderPrivateKey			*string `json:"sender_private_key"`
	SenderBlockchainAddress		*string `json:"sender_blockchain_address"`
	ReceiverBlockchainAddress	*string `json:"receiver_blockchain_address"`
	SenderPublicKey 			*string `json:"sender_public_key"`
	Value						*string `json:"value"`
}

func (tr *TransactionRequest) Validate() (bool, string) {
	if tr.SenderPrivateKey == nil {
		return false, "Sender private key is missing"
	}
	if tr.SenderBlockchainAddress == nil {
		return false, "Sender blockchain address is missing"
	}
	if tr.ReceiverBlockchainAddress == nil {
		return false, "Receiver blockchain address is missing"
	}
	if tr.SenderPublicKey == nil {
		return false, "Sender public key is missing"
	}
	if tr.Value == nil {
		return false, "Transaction value is missing"
	}
	return true, "" // Validation succeeded, no error message
}
