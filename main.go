package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Block struct {
	Pos       int
	Data      SneakerCheckout
	TimeStamp string
	Hash      string
	PrevHash  string
}

type SneakerCheckout struct {
	SneakerID    string `json:"sneaker_id"`
	Customer     string `json:"customer"`
	CheckoutDate string `json:"checkout_date"`
	IsGenesis    bool   `json:"is_genesis"` // Genesis is the intial block in a blockchain
}

type Sneaker struct {
	ID              string `json:"id"`
	Silhouette      string `json:"silhouette"`
	Brand           string `json:"brand"`
	ManufactureDate string `json:"manufacture_date"`
	ISBN            string `json:"isbn"`
}

type Blockchain struct {
	blocks []*Block
}

var Blockchain *Blockchain

func main() {
	fmt.Println("BlockKicks")
	r := mux.NewRouter()
	r.HandleFunc("/", getBlockchain).Methods("GET")
	r.HandleFunc("/", writeBlock).Methods("POST")
	r.HandleFunc("/new", newSneaker).Methods("POST")

	log.Fatal(http.ListenAndServe(":4000", r))
	fmt.Println("Listening at port 4000...")
}
