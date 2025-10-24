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
}

type BlockData struct {
	Sender Address
	Receiver Address
	Amount float64
}

