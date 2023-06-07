package main 

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"log"
	"os"
)

type Block struct {
	data            map[string]interface{}
	prevHash, hash  string
	time            time.Time
	pow             int
}

type Blockchain struct {
	genesis         Block 
	chain           []Block
	difficulty      int
}

func (block Block) calcHash() string {
	data, _  := json.Marshal(block.data)
	blockData := block.prevHash + string(data) + block.time.String() + strconv.Itoa(block.pow)
	blockHash := sha256.Sum256([]byte(blockData))
	return string(blockHash[:])
}

func (block *Block) mine(difficulty int) {
	neededPrefix := strings.Repeat("0", difficulty)
	for !strings.HasPrefix(block.hash, neededPrefix) {
		block.pow ++ 

		// log.Println("mining hashes...")
		block.hash = block.calcHash()
	}
}

func Create(difficulty int) Blockchain {
	genesis := Block {
		hash: "0",
		time: time.Now(),
	}

	newBlockChain := Blockchain {
		genesis: genesis,
		chain: []Block{genesis},
		difficulty: difficulty,
	}

	log.Println("created a new blockchain.")
	return newBlockChain
}

func (blockchain *Blockchain) addBlock(sender, receiver string, amount float64) {
	newData := map[string]interface{} {
		"Sender": sender,
		"Receiver": receiver,
		"Amount": amount,
	}

	newBlock := Block {
		data: newData,
		prevHash: blockchain.chain[len(blockchain.chain) - 1].hash,
		time: time.Now(),
	}

	newBlock.mine(blockchain.difficulty)

	blockchain.chain = append(blockchain.chain, newBlock)
	log.Println("added a new block.")
}

func (blockchain Blockchain) validator() bool {
	for i := 2; i < len(blockchain.chain); i++ {
		if blockchain.chain[i - 1].hash != blockchain.chain[i - 1].calcHash() ||
		blockchain.chain[i].hash != blockchain.chain[i].calcHash() {
			log.Println("failed to validate blockchain.")
			return false
		}
	}

	log.Println("blockchain is valid.")
	return true
}

func (blockchain Blockchain) display() {
	for i := 1; i < len(blockchain.chain); i++ {
		cur := blockchain.chain[i]

		fmt.Println("----------------------------- Block #" + strconv.Itoa(i) + " -----------------------------")
		fmt.Println("Sender:", cur.data["Sender"])
		fmt.Println("Receiver:", cur.data["Receiver"])
		fmt.Println("Amount:", cur.data["Amount"])
		// fmt.Printf("Previous block's hash value: %x\n", cur.prevHash)
		// fmt.Printf("Current block's hash value: %x\n", cur.hash)
		fmt.Println("Time created:", cur.time)
	}
}
func main() {
    file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)

    if err != nil {
        log.Fatal(err)
    }

    log.SetOutput(file)

	blockchain := Create(1)

    blockchain.addBlock("Alice", "Bob", 5)
    blockchain.addBlock("John", "Bob", 2)
	blockchain.addBlock("Bob", "Alice", 9.99)

    if blockchain.validator() {
		blockchain.display()
	}

	log.Println("session ended.")
}
