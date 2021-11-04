package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ParkChangSun/nomadcoin/blockchain"
	"github.com/ParkChangSun/nomadcoin/utils"
	"github.com/gorilla/mux"
)

var port string

type url string

func (u url) MarshalText() (text []byte, err error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

type urlDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

type balanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{"/", "GET", "See documentation", ""},
		{"/blocks", "POST", "Add a block", "data:string"},
		{"/blocks", "GET", "All blocks", ""},
		{"/status", "GET", "see status", ""},
	}
	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(rw).Encode(blockchain.Blockchain().Blocks())
	case "POST":
		blockchain.Blockchain().AddBlock()
		rw.WriteHeader(http.StatusCreated)
	}
}

type errorResponse struct {
	Message string `json:"message"`
}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	block, err := blockchain.FindBlock(hash)
	encoder := json.NewEncoder(rw)
	if err == blockchain.ErrNotFound {
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	} else {
		encoder.Encode(block)

	}
}

func status(rw http.ResponseWriter, r *http.Request) {
	json.NewEncoder(rw).Encode(blockchain.Blockchain())
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	total := r.URL.Query().Get("total")
	switch total {
	case "true":
		res := balanceResponse{address, blockchain.Blockchain().BalanceByAddress(address)}
		utils.HandleError(json.NewEncoder(rw).Encode(res))
	default:
		utils.HandleError(json.NewEncoder(rw).Encode(blockchain.Blockchain().UTxOutsByAddress(address)))

	}
}

func mempool(rw http.ResponseWriter, r *http.Request) {
	utils.HandleError(json.NewEncoder(rw).Encode(blockchain.Mempool.Txs))
}

type addTxPayload struct {
	To     string
	Amount int
}

func transactions(rw http.ResponseWriter, r *http.Request) {
	var payload addTxPayload
	utils.HandleError(json.NewDecoder(r.Body).Decode(&payload))
	err := blockchain.Mempool.AddTx(payload.To, payload.Amount)
	if err != nil {
		json.NewEncoder(rw).Encode(errorResponse{err.Error()})
	}
	rw.WriteHeader(http.StatusCreated)
}

func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)
	router := mux.NewRouter()
	router.Use(jsonContentTypeMiddleware)
	router.HandleFunc("/", documentation)
	router.HandleFunc("/blocks", blocks)
	router.HandleFunc("/block/{hash:[a-f0-9]+}", block)
	router.HandleFunc("/status", status)

	router.HandleFunc("/balance/{address}", balance)
	router.HandleFunc("/mempool", mempool)
	router.HandleFunc("/transactions", transactions)
	log.Fatal(http.ListenAndServe(port, router))
}
