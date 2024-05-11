package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
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

func newSneaker(w http.ResponseWriter, r *http.Request) {
	var sneaker Sneaker

	if err := json.NewDecoder(r.Body).Decode(&sneaker); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Could not create:%v", err)
		w.Write([]byte("Could not add new sneaker"))
		return
	}

	h := md5.New()
	io.WriteString(h, sneaker.ISBN+sneaker.ManufactureDate)
	sneaker.ID = fmt.Sprintf("%x", h.Sum(nil))

	resp, err := json.MarshalIndent(sneaker, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Could not marshal payload: %v", err)
		w.Write([]byte("Could not save book data"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
func main() {
	fmt.Println("BlockKicks")
	r := mux.NewRouter()
	r.HandleFunc("/", getBlockchain).Methods("GET")
	r.HandleFunc("/", writeBlock).Methods("POST")
	r.HandleFunc("/new", newSneaker).Methods("POST")

	log.Fatal(http.ListenAndServe(":4000", r))
	fmt.Println("Listening at port 4000...")
}
