package blockchain

import (
	"fmt"
	"time"
)

type Block struct {
	Index int
	Timestamp time.Time
	Data BlockData
	Hash string
	PrevHash string
	nonce int
}

type BlockData struct {
	Sender Address
	Receiver Address
	Amount float64
}

