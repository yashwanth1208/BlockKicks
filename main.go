package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

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
	ArticleName     string `json:"article_name"`
}

type Blockchain struct {
	blocks []*Block
}

var BlockChain *Blockchain

func (b *Block) generateHash() {

	bytes, _ := json.Marshal(b.Data)

	data := string(b.Pos) + b.TimeStamp + string(bytes) + b.PrevHash

	hash := sha256.New()
	hash.Write([]byte(data))
	b.Hash = hex.EncodeToString(hash.Sum(nil))
}

func CreateBlock(prevBlock *Block, checkoutItem SneakerCheckout) *Block {
	block := &Block{}
	block.Pos = prevBlock.Pos + 1
	block.TimeStamp = time.Now().String()
	block.Data = checkoutItem
	block.PrevHash = prevBlock.Hash
	block.generateHash()

	return block
}

// Struct Method
func (bc *Blockchain) AddBlock(data SneakerCheckout) {

	prevBlock := bc.blocks[len(bc.blocks)-1]

	block := CreateBlock(prevBlock, data)

	if validBlock(block, prevBlock) {
		bc.blocks = append(bc.blocks, block)
	}
}

func validBlock(block, prevBlock *Block) bool {

	if prevBlock.Hash != block.PrevHash {
		return false
	}

	if !block.validateHash(block.Hash) {
		return false
	}

	if prevBlock.Pos+1 != block.Pos {
		return false
	}

	return true
}

func (b *Block) validateHash(hash string) bool {
	b.generateHash()
	if b.Hash != hash {
		return false
	}
	return true
}
func writeBlock(w http.ResponseWriter, r *http.Request) {
	var checkoutItem SneakerCheckout

	if err := json.NewDecoder(r.Body).Decode(&checkoutItem); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Could not write block:%v", err)
		w.Write([]byte("Could not write block"))
	}

	BlockChain.AddBlock(checkoutItem)
	resp, err := json.MarshalIndent(checkoutItem, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not marshal payload: %v", err)
		w.Write([]byte("could not write block"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func newSneaker(w http.ResponseWriter, r *http.Request) {
	var sneaker Sneaker

	if err := json.NewDecoder(r.Body).Decode(&sneaker); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Could not create:%v", err)
		w.Write([]byte("Could not add new sneaker"))
		return
	}

	h := md5.New()
	io.WriteString(h, sneaker.ArticleName+sneaker.ManufactureDate)
	sneaker.ID = fmt.Sprintf("%x", h.Sum(nil))

	resp, err := json.MarshalIndent(sneaker, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Could not marshal payload: %v", err)
		w.Write([]byte("Could not save sneaker data"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
func GenesisBlock() *Block {
	return CreateBlock(&Block{}, SneakerCheckout{IsGenesis: true})
}
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{GenesisBlock()}}
}

func getBlockchain(w http.ResponseWriter, r *http.Request) {
	jbytes, err := json.MarshalIndent(BlockChain.blocks, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	io.WriteString(w, string(jbytes))
}
func main() {

	BlockChain = NewBlockchain()
	fmt.Println("BlockKicks")
	r := mux.NewRouter()
	r.HandleFunc("/", getBlockchain).Methods("GET")
	r.HandleFunc("/", writeBlock).Methods("POST")
	r.HandleFunc("/new", newSneaker).Methods("POST")

	go func() {
		for _, block := range BlockChain.blocks {
			fmt.Printf("Prev Hash: %x\n", block.PrevHash)
			bytes, _ := json.MarshalIndent(block.Data, "", " ")
			fmt.Printf("Data:%v\n", string(bytes))
			fmt.Printf("Hash:%x\n", block.Hash)
			fmt.Println()
		}

	}()
	log.Fatal(http.ListenAndServe(":4000", r))
	fmt.Println("Listening at port 4000...")
}
