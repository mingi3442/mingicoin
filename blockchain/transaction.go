package blockchain

import (
	"errors"
	"sync"
	"time"

	"github.com/mingi3442/mingicoin/utils"
	"github.com/mingi3442/mingicoin/wallet"
)

const (
	minerReward int = 50
)

type mempool struct {
	Txs map[string]*Tx
	m   sync.Mutex
}

var m *mempool //empty mempool 생성
var memOnce sync.Once

func Mempool() *mempool {
	memOnce.Do(func() {
		m = &mempool{
			Txs: make(map[string]*Tx),
		}
	})
	return m
}

type Tx struct {
	ID        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

type TxIn struct {
	TxID      string `json:"txId"`
	Index     int    `json:"index"`
	Signature string `json:"signature"`
}
type TxOut struct {
	Address string `json:"address"`
	Amount  int    `json:"amount"`
}
type UTxOut struct {
	TxID   string `json:"txId"`
	Index  int    `json:"index"`
	Amount int    `json:"amount"`
}

func (t *Tx) getId() {
	t.ID = utils.Hash(t)
}

func (t *Tx) sign() {
	//txIn을 생성시, privateKey를 이용해서 sign을 한다
	for _, txIn := range t.TxIns {
		txIn.Signature = wallet.Sign(t.ID, wallet.Wallet())
	}
}

func validate(tx *Tx) bool { //검증할 transaction을 가져온다
	valid := true
	for _, txIn := range tx.TxIns {
		prevTx := FindTx(Blockchain(), txIn.TxID)
		if prevTx == nil {
			valid = false
			break
		}
		address := prevTx.TxOuts[txIn.Index].Address //가져온 Tx를 참조하는 TxOut을 찾아 Address를 찾음
		valid = wallet.Verify(txIn.Signature, tx.ID, address)
		if !valid {
			break
		}
	}
	return valid
}

func isOnMempool(UTxOut *UTxOut) bool {
	exists := false
Outer: //label
	for _, tx := range Mempool().Txs {
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

var ErrorNoMoney = errors.New("Not enough Money..!")
var ErrorNotValid = errors.New("Tx Invalid..")

func makeTx(from, to string, amount int) (*Tx, error) {
	if BalanceByAddress(from, Blockchain()) < amount {
		return nil, ErrorNoMoney
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
	tx.sign()
	valid := validate(tx)
	if !valid {
		return nil, ErrorNotValid
	}
	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) (*Tx, error) {
	tx, err := makeTx(wallet.Wallet().Address, to, amount)
	if err != nil {
		return nil, err
	}
	m.Txs[tx.ID] = tx
	return tx, nil
}

func (m *mempool) TxToConfirm() []*Tx {
	coinbase := makeCoinbaseTx(wallet.Wallet().Address) //coinbase의 transaction
	var txs []*Tx                                       //mempool에 있는 모든 transaction
	txs = append(txs, coinbase)                         //coinbaseTx + mempoolTx
	for _, tx := range m.Txs {
		txs = append(txs, tx)
	}
	m.Txs = make(map[string]*Tx) //mempool을 비움
	return txs
}

func (m *mempool) AddPeerTx(tx *Tx) {
	m.m.Lock()
	defer m.m.Unlock()
	m.Txs[tx.ID] = tx
}
