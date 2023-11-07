package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
	"text/template"

	"github.com/palmcivet7/go-blockchain/block"
	"github.com/palmcivet7/go-blockchain/utils"
	"github.com/palmcivet7/go-blockchain/wallet"
)

const tempDir = "templates"

type WalletServer struct {
	port		uint16
	gateway		string
}

func NewWalletServer(port uint16, gateway string) *WalletServer {
	return &WalletServer{port, gateway}
}

func (ws *WalletServer) Port() uint16 {
	return ws.port
}

func (ws *WalletServer) Gateway() string {
	return ws.gateway
}

func (ws *WalletServer) Index(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		t, err := template.ParseFiles(path.Join(tempDir, "index.html"))
		if err != nil {
			log.Printf("Error parsing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if err := t.Execute(w, ""); err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	default:
		log.Printf("ERROR: Invalid HTTP Method")
	}
}

func (ws *WalletServer) Wallet(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		w.Header().Add("Content-Type", "application-json")
		myWallet := wallet.NewWallet()
		m, _ := myWallet.MarshalJSON()
		io.WriteString(w, string(m[:]))
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")
	}
}

func (ws *WalletServer) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	switch r.Method{
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var t wallet.TransactionRequest
		err := decoder.Decode(&t)
		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}
		isValid, errorMsg := t.Validate()
		if !isValid {
			log.Printf("ERROR: %s", errorMsg)
			io.WriteString(w, string(utils.JsonStatus(errorMsg)))
			return
		}

		publicKey := utils.PublicKeyFromString(*t.SenderPublicKey)
		privateKey := utils.PrivateKeyFromString(*t.SenderPrivateKey, publicKey)
		value, err := strconv.ParseFloat(*t.Value, 32)
		if err != nil {
			log.Println("ERROR: parse error ")
			io.WriteString(w, string(utils.JsonStatus(errorMsg)))
			return
		}
		value64 := float64(value)

		w.Header().Add("Content-Type", "application-json")

		transaction := wallet.NewTransaction(privateKey, publicKey,
			*t.SenderBlockchainAddress, *t.ReceiverBlockchainAddress, value64)
		signature := transaction.GenerateSignature()
		signatureStr := signature.String()

		bt := &block.TransactionRequest{
			t.SenderBlockchainAddress,
			t.ReceiverBlockchainAddress,
			t.SenderPublicKey,
			&value64, &signatureStr,
		}
		m, _ := json.Marshal(bt)
		buf := bytes.NewBuffer(m)

		resp, _ := http.Post(ws.Gateway() + "/transactions", "application/json", buf)
		if resp.StatusCode == 201 {
			io.WriteString(w, string(utils.JsonStatus("success")))
			return
		}
		io.WriteString(w, string(utils.JsonStatus("fail ")))
		 
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")
	}
}


func (ws *WalletServer) Run() {
	http.HandleFunc("/", ws.Index)
	http.HandleFunc("/wallet", ws.Wallet)
	http.HandleFunc("/transaction", ws.CreateTransaction)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(ws.Port())), nil))
}