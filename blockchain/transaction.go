package blockchain

import (
	"errors"
	"time"

	"github.com/ParkChangSun/nomadcoin/utils"
)

type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{}

type Tx struct {
	Id        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txins"`
	TxOuts    []*TxOut `json:"txouts"`
}

type TxIn struct {
	TxId  string
	Index int
	Owner string `json:"owner"`
}

type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

type UTxOut struct {
	TxId   string
	Index  int
	Amount int
}

func (t *Tx) getId() {
	t.Id = utils.Hash(t)
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"", -1, "coinbase"},
	}

	txOuts := []*TxOut{
		{address, 50},
	}

	tx := Tx{
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return &tx
}

func isOnMempool(uTx *UTxOut) (exists bool) {
	for _, tx := range Mempool.Txs {
		for _, txIn := range tx.TxIns {
			exists = uTx.Index == txIn.Index && uTx.TxId == txIn.TxId
		}
	}
	return
}

func makeTx(from, to string, amount int) (*Tx, error) {
	if Blockchain().BalanceByAddress(from) < amount {
		return nil, errors.New("not enough money")
	}
	var txOuts []*TxOut
	var txIns []*TxIn

	uTxOuts := Blockchain().UTxOutsByAddress(from)
	total := 0
	for _, u := range uTxOuts {
		if total >= amount {
			break
		}
		txIns = append(txIns, &TxIn{u.TxId, u.Index, from})
		total += u.Amount
	}
	if change := total - amount; change > 0 {
		txOuts = append(txOuts, &TxOut{from, change})
	}
	txOuts = append(txOuts, &TxOut{to, amount})
	tx := &Tx{
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx("park", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func (m *mempool) TxToConfirm() []*Tx {
	coinbase := makeCoinbaseTx("park")
	txs := m.Txs
	txs = append(txs, coinbase)
	m.Txs = nil
	return txs
}
