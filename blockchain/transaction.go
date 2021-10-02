package blockchain

import (
	"errors"
	"time"

	"github.com/mingi3442/mingicoin/utils"
)

const (
	minerReward int = 50
)

type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{} //empty mempool 생성

type Tx struct {
	ID        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

func (t *Tx) getId() {
	t.ID = utils.Hash(t)
}

type TxIn struct {
	TxID  string `json:"txId"`
	Index int    `json:"index"`
	Owner string `json:"owner"`
}
type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}
type UTxOut struct {
	TxID   string `json:"txId"`
	Index  int    `json:"index"`
	Amount int    `json:"amount"`
}

func isOnMempool(UTxOut *UTxOut) bool {
	exists := false
Outer: //label
	for _, tx := range Mempool.Txs {
		for _, input := range tx.TxIns {
			exists = input.TxID == UTxOut.TxID && input.Index == UTxOut.Index
			break Outer //for문이 중첩되었을 때 label을 이용하면 조건이 만족했을때 label로 빠져나간다.
		}
	}
	return exists
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"", -1, "COINBASE"},
	}
	txOuts := []*TxOut{
		{address, minerReward},
	}
	tx := Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return &tx
}

func makeTx(from, to string, amount int) (*Tx, error) {
	if BalanceByAddress(from, Blockchain()) < amount {
		return nil, errors.New("Not enough Money..!")
	}
	var txOuts []*TxOut
	var txIns []*TxIn
	total := 0
	uTxOuts := UTxOutsByAddress(from, Blockchain())
	for _, uTxOut := range uTxOuts {
		if total >= amount {
			break
		}
		txIn := &TxIn{uTxOut.TxID, uTxOut.Index, from}
		txIns = append(txIns, txIn)
		total += uTxOut.Amount
	}
	if change := total - amount; change != 0 {
		changeTxOut := &TxOut{from, change}
		txOuts = append(txOuts, changeTxOut)
	}
	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)
	tx := &Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx("mingi", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func (m *mempool) TxToConfirm() []*Tx {
	coinbase := makeCoinbaseTx("mingi") //coinbase의 transaction
	txs := m.Txs                        //mempool에 있는 모든 transaction
	txs = append(txs, coinbase)         //coinbaseTx + mempoolTx
	m.Txs = nil                         //mempool을 비움
	return txs
}
