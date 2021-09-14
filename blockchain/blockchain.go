package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type Block struct {
	Data     string
	Hash     string
	PrevHash string
}
type blockchain struct {
	blocks []*Block
}

var b *blockchain
var once sync.Once

func (b *Block) getHash() {
	hash := sha256.Sum256([]byte(b.Data + b.PrevHash))
	b.Hash = fmt.Sprintf("%x", hash)
}

func getLastHash() string {
	totalBlock := len(GetBlockchain().blocks)
	if totalBlock == 0 {
		return ""
	}
	return GetBlockchain().blocks[totalBlock-1].Hash
}

func createBlock(data string) *Block {
	newBlock := Block{data, "", getLastHash()}
	newBlock.getHash()
	return &newBlock
}

func (b *blockchain) AddBlock(data string) {
	b.blocks = append(b.blocks, createBlock(data))
}

func GetBlockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{}
			b.AddBlock("Genesis")
		})
	}
	return b
}
func (b *blockchain) AllBlock() []*Block {
	return b.blocks
}
