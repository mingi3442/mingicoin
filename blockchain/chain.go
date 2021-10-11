package blockchain

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/mingi3442/mingicoin/db"
	"github.com/mingi3442/mingicoin/utils"
)

const (
	defaultDifficulty  int = 2
	difficultyInterval int = 5
	blockInterval      int = 2
	allowedRange       int = 2
)

type blockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentDifficulty"`
	m                 sync.Mutex
}

var b *blockchain
var once sync.Once

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) AddBlock() *Block {
	block := createBlock(b.NewestHash, b.Height+1, getDifficulty(b))
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	persistBlockchain(b)
	return block
}
func persistBlockchain(b *blockchain) {
	db.SaveCheckpoint(utils.ToBytes(b))
}

func Blocks(b *blockchain) []*Block {
	b.m.Lock()
	defer b.m.Unlock()
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

func Txs(b *blockchain) []*Tx {
	var txs []*Tx
	for _, block := range Blocks(b) {
		txs = append(txs, block.Transactions...)
	}
	return txs
}

func FindTx(b *blockchain, targetID string) *Tx {
	for _, tx := range Txs(b) {
		if tx.ID == targetID {
			return tx
		}
	}
	return nil
}

func recalculateDifficulty(b *blockchain) int {
	allBlocks := Blocks(b)
	newestBlock := allBlocks[0]
	lastReacalculatedBlock := allBlocks[difficultyInterval-1]
	actualTime := (newestBlock.Timestamp / 60) - (lastReacalculatedBlock.Timestamp / 60)
	expectedTime := difficultyInterval * blockInterval
	if actualTime <= (expectedTime - allowedRange) {
		return b.CurrentDifficulty + 1
	} else if actualTime >= (expectedTime + allowedRange) {
		return b.CurrentDifficulty - 1
	}
	return b.CurrentDifficulty

}

func getDifficulty(b *blockchain) int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		return recalculateDifficulty(b)
	} else {
		return b.CurrentDifficulty
	}
}

func UTxOutsByAddress(addres string, b *blockchain) []*UTxOut {
	var uTxOuts []*UTxOut
	creatorTxs := make(map[string]bool)
	for _, block := range Blocks(b) {
		for _, tx := range block.Transactions {
			for _, input := range tx.TxIns {
				if input.Signature == "COINBASE" {
					break
				}
				if FindTx(b, input.TxID).TxOuts[input.Index].Address == addres {
					creatorTxs[input.TxID] = true
				}
			}
			for index, output := range tx.TxOuts {
				if output.Address == addres {
					if _, ok := creatorTxs[tx.ID]; !ok {
						uTxOut := &UTxOut{tx.ID, index, output.Amount}
						if !isOnMempool(uTxOut) {
							uTxOuts = append(uTxOuts, uTxOut)

						}
					}
				}
			}
		}
	}
	return uTxOuts
}

func BalanceByAddress(address string, b *blockchain) int {
	txOuts := UTxOutsByAddress(address, b) //b.TxOutsByAddress를 이용해서 특정 address의 txOuts들 가져오기
	var amount int                         //변수 선언
	for _, txOut := range txOuts {
		amount += txOut.Amount //for loop를 이용해서 txOuts안의 Amount를 선언한 변수에 더한다
	}
	return amount
}

func Blockchain() *blockchain {
	once.Do(func() {
		b = &blockchain{Height: 0}    //empty blockchain을 만듦
		checkpoint := db.Checkpoint() //db에서 blockchain checkpoint를 찾는다
		if checkpoint == nil {
			b.AddBlock() //checkpoint가 nil일 경우 blockchain을 initalize
		} else {
			b.restore(checkpoint) //찾을 경우 checkpoint로 부터 blockchain 복원
		}
	})
	return b
}

func Status(b *blockchain, rw http.ResponseWriter) {
	b.m.Lock()
	defer b.m.Unlock()
	json.NewEncoder(rw).Encode(b)
}

func (b *blockchain) Repalce(newBlocks []*Block) {
	//blockchain update
	b.m.Lock()
	defer b.m.Unlock()
	b.CurrentDifficulty = newBlocks[0].Difficulty
	b.Height = len(newBlocks)
	b.NewestHash = newBlocks[0].Hash
	persistBlockchain(b)
	db.EmptyBlocks()
	for _, block := range newBlocks {
		persistBlock(block)
	}
}

func (b *blockchain) AddPeerBlock(newBlock *Block) {
	b.m.Lock()
	m.m.Lock()
	defer b.m.Unlock()
	defer m.m.Unlock()
	b.Height += 1
	b.CurrentDifficulty = newBlock.Difficulty
	b.NewestHash = newBlock.Hash

	persistBlockchain(b)
	persistBlock(newBlock)

	for _, tx := range newBlock.Transactions {
		_, ok := m.Txs[tx.ID]
		if ok {
			delete(m.Txs, tx.ID)
		}
	}

}
