package blockchain

import (
	"fmt"
	"time"
	"encoding/json"
	"crypto/sha256"
	"github.com/oliinant/fish-ham/scripts"
)


type Transaction struct {
	ID int `json:"id"`
	Sender string `json:"sender"`
	Receiver string `json:"receiver"`
	Amount float64 `json:"amount"`
	Tax float64 `json:"tax"`
}

type TxList struct {
	Transactions []Transaction `json:"transactions"`
}

func (txs *TxList) getTxListInfoMap() []map[string]interface{} {
	txListMap := []map[string]interface{}{}

	for _, tx := range txs.Transactions {
		txListMap = append(txListMap, scripts.GetInfoMap(tx))
	}
	return txListMap
}

type Hash string

func HashIt(s string) (Hash, error) {
	if s == "" {
		return "", fmt.Errorf("Failed hash conversion: Empty string")
	}

	for _, r  := range s {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F' )) {
			return "", fmt.Errorf("Invalid character in hash: %q", r)
		}
	}
	return Hash(s), nil
}

type Block struct {
	Index int `json:"index"`
	Timestamp time.Time `json:"timestamp"`
	Transactions TxList `json:"transactions"`
	Hash Hash `json:"hash"`
	PrevHash Hash `json:"prev_hash"`
	Nonce int `json: "nonce"`
	Reward float64 `json:"reward"`
}

func (b *Block) Serializer() (string, error) {
	delimiter := "\\fh~8"

	dataJSON, err := json.Marshal(b.Transactions)
	if err != nil {
		return "", scripts.WrapError("Failed to serialize block", err)
	}

	serializedBlock := fmt.Sprint(
		b.Index, delimiter,
		b.Timestamp, delimiter,
		string(dataJSON), delimiter,
		b.PrevHash, delimiter,
		b.Nonce,
	)
	return serializedBlock, nil
}

func (b *Block) CalculateHash() (Hash, error) {
	errorMessage := "Failed to calculate hash"

	serializedBlock, errSerialze := b.Serializer()
	if errSerialze != nil {
		return "", scripts.WrapError(errorMessage, errSerialze)
	}

	byteHash := sha256.Sum256([]byte(serializedBlock))
	hexHash := fmt.Sprintf("%x", byteHash)

	hashHash, errHash := HashIt(hexHash)
	if errHash != nil {
		return "", scripts.WrapError(errorMessage, errHash)
	}

	return hashHash, nil
}

type Blockchain struct {
	Chain map[Hash]*Block
	Tips map[Hash]*Block
	Genesis *Block 
}