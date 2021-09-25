package blockchain

import (
	"sync"

	"github.com/mingi3442/mingicoin/db"
	"github.com/mingi3442/mingicoin/utils"
)

type blockchain struct {
	NewestHash string `json:"newestHash"`
	Height     int    `json:"height"`
}

var b *blockchain
var once sync.Once

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) persist() {
	db.SaveBlockchain(utils.ToBytes(b))
}

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.persist()
}

func (b *blockchain) Blocks() []*Block {
	var blocks []*Block        // blocks slice를 만듦
	hashCursor := b.NewestHash //처음 찾는 hash
	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break //block의 PrevHash가 ""이면, 다시말해 Genesis Block이면 for loop를 중단하도 blocks return
		}
	}
	return blocks
}

func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{"", 0}        //empty blockchain을 만듦
			checkpoint := db.Checkpoint() //db에서 blockchain checkpoint를 찾는다
			if checkpoint == nil {
				b.AddBlock("Genesis") //checkpoint가 nil일 경우 blockchain을 initalize
			} else {
				b.restore(checkpoint) //찾을 경우 checkpoint로 부터 blockchain 복원
			}
		})
	}
	return b
}
