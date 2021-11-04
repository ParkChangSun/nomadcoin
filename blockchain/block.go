package blockchain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ParkChangSun/nomadcoin/db"
	"github.com/ParkChangSun/nomadcoin/utils"
)

type Block struct {
	Hash         string `json:"hash"`
	PrevHash     string `json:"prevhash,omitempty"`
	Height       int    `json:"height"`
	Difficulty   int    `json:"difficulty"`
	Nonce        int    `json:"nonce"`
	Timestamp    int    `json:"timestamp"`
	Transactions []*Tx  `json:"transactions"`
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

var ErrNotFound = errors.New("not found")

func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

func FindBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restore(blockBytes)
	return block, nil
}

func createBlock(prevHash string, height int) *Block {
	block := Block{
		PrevHash:   prevHash,
		Height:     height,
		Difficulty: Blockchain().difficulty(),
	}
	block.mine()
	block.Transactions = Mempool.TxToConfirm()
	block.persist()
	return &block
}

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		b.Timestamp = int(time.Now().Unix())
		hash := utils.Hash(b)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			fmt.Printf("\n%d\n%s\n", b.Difficulty, b.Hash)
			return
		} else {
			b.Nonce++
		}
	}
}
