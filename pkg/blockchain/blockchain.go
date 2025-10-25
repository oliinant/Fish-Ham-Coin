package blockchain

import (
	"fmt"
	"time"
	"encoding/json"
	"crypto/sha256"
	"unicode"
	"fish-ham/scripts"
)


type Transaction struct {
	Sender Address
	Receiver Address
	Amount float64
}

type Hash string

func HashIt(s string) (Hash, error) {
	for _, r  := range s {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F' )) {
			return "", fmt.Errorf("Invalid character in hash: %q", r)
		}
	}
	return Hash(s), nil
}

type Block struct {
	Index int
	Timestamp time.Time
	Data []Transaction
	Hash Hash
	PrevHash Hash
	Nonce int 
}

func (b *Block) Serializer() (string, error) {
	delimiter := "\\fh~8"

	dataJSON, err := json.Marshal(b.Data)
	if err != nil {
		return "", scripts.WrapError("Failed to serialize block", err)
	}

	serializedBlock := fmt.Sprint(
		b.Index, delimiter,
		b.Timestamp, delimiter,
		string(dataJSON), delimiter,
		b.PrevHash, delimiter,
		b.Nonce
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