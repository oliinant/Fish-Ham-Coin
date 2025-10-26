package blockchain

import (
	"fmt"
	"time"
	"encoding/json"
	"crypto/sha256"
	"github.com/oliinant/fish-ham/scripts"
	"github.com/google/uuid"
)


type Transaction struct {
	ID uuid.UUID `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Sender string `json:"sender"`
	Receiver string `json:"receiver"`
	Amount float64 `json:"amount"`
	Tax float64 `json:"tax"`
}

func CustomTxInit(time time.Time, sender string, receiver string, amount float64, taxPercent float64) *Transaction {
	taxDecimal := taxPercent / 100

	return &Transaction{
		ID: uuid.New(),
		Timestamp: time,
		Sender: sender,
		Receiver: receiver,
		Amount: amount,
		Tax: amount * taxDecimal,
	}
}

// Default means time is when transaction created and tax is 1%
func DefaultTxInit(sender string, receiver string, amount float64) *Transaction {
	return &Transaction{
		ID: uuid.New(),
		Timestamp: time.Now(),
		Sender: sender,
		Receiver: receiver,
		Amount: amount,
		Tax: amount * 0.01,
	}
}

type TxList struct {
	Transactions map[uuid.UUID]*Transaction `json:"transactions"`
}

func (txs *TxList) TxListInfoMap() []map[string]interface{} {
	txListMap := []map[string]interface{}{}

	for _, tx := range txs.Transactions {
		info := scripts.InfoMap(tx)
		txListMap = append(txListMap, info)
	}
	return txListMap
}

func (txs *TxList) AddTransaction(tx *Transaction) {
	if txs.Transactions == nil {
		txs.Transactions = make(map[uuid.UUID]*Transaction)
	}
	txs.Transactions[tx.ID] = tx
}

func (txs *TxList) RemoveTransactionByID(id uuid.UUID) {
	delete(txs.Transactions, id)
}

func (txs *TxList) TransactionByID(id uuid.UUID) (*Transaction, error) {
	tx, ok := txs.Transactions[id]
	if !ok {
		return nil, fmt.Errorf("Transaction %s not found", id.String())
	}
	return tx, nil
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
	Nonce int `json:"nonce"`
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