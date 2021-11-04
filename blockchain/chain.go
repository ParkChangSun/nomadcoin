package blockchain

import (
	"sync"

	"github.com/ParkChangSun/nomadcoin/db"
	"github.com/ParkChangSun/nomadcoin/utils"
)

type blockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentdifficulty"`
}

const defaultDifficulty = 2
const difficultyInterval = 5
const blockInterval = 2
const allowedRange = 2

var b *blockchain
var once sync.Once

func (b *blockchain) persist() {
	db.SaveBlockchain(utils.ToBytes(b))
}

func (b *blockchain) AddBlock() {
	block := createBlock(b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	b.persist()
}

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{Height: 0}
			checkpoint := db.Checkpoint()
			if checkpoint == nil {
				b.AddBlock()
			} else {
				b.restore(checkpoint)
			}
		})
	}
	return b
}

func (b *blockchain) Blocks() (blocks []*Block) {
	hashCursor := b.NewestHash
	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return
}

func (b *blockchain) difficulty() int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		return b.recalcDiff()
	} else {
		return b.CurrentDifficulty
	}
}

func (b *blockchain) recalcDiff() int {
	allBlocks := b.Blocks()
	newestBlock := allBlocks[0]
	lastRecalculatedBlock := allBlocks[blockInterval-1]
	actualTime := (newestBlock.Timestamp - lastRecalculatedBlock.Timestamp) / 60
	expectedTime := blockInterval * difficultyInterval
	if actualTime < expectedTime-allowedRange {
		return b.CurrentDifficulty + 1
	} else if actualTime > expectedTime+allowedRange {
		return b.CurrentDifficulty - 1
	} else {
		return b.CurrentDifficulty
	}
}

func (b *blockchain) UTxOutsByAddress(address string) (ownedTxOuts []*UTxOut) {
	usedTxs := make(map[string]bool)

	for _, block := range b.Blocks() {
		for _, tx := range block.Transactions {
			for _, input := range tx.TxIns {
				if input.Owner == address {
					usedTxs[input.TxId] = true
				}
			}
			for index, output := range tx.TxOuts {
				if output.Owner == address {
					if _, ok := usedTxs[tx.Id]; !ok {
						uTxOut := &UTxOut{tx.Id, index, output.Amount}
						if !isOnMempool(uTxOut) {
							ownedTxOuts = append(ownedTxOuts, uTxOut)
						}
					}
				}
			}
		}
	}
	return
}

func (b *blockchain) BalanceByAddress(address string) (balance int) {
	txOuts := b.UTxOutsByAddress(address)
	for _, v := range txOuts {
		balance += v.Amount
	}
	return
}
