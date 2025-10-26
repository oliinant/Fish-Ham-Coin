package blockchain

import (
	"testing"
	"github.com/oliinant/fish-ham/scripts"
	"time"
	"encoding/json"
	"fmt"
)
// Helpers

func newTestBlock() *Block {
	return &Block {
		Index: 1,
		Timestamp: time.Date(2025, time.October, 23, 15, 0, 56, 0, time.UTC),
		Transactions: TxList{
			Transactions: []Transaction{
			    {ID: 1, Sender: "Saul", Receiver: "Bush", Amount: 10, Tax: 0.1},
		}},
		PrevHash: "0000",
		Nonce: 69,
	}
}

// Unit Tests

func TestHashIt(t *testing.T) {
	tests := []scripts.TestCase[string, Hash]{
		{"valid lowercase hex", "abc123", Hash("abc123"), false},
		{"valid uppercase hex", "ABC123", Hash("ABC123"), false},
		{"invalid letter hex", "xyz123", "", true},
		{"empty string", "", "", true},
	}

	scripts.BoilerTestFunc[string, Hash](t, HashIt, tests)
}

//func TestgetTxListInfoMap(t *testing.T) {
	//tests := []scripts.TestCase[]{

	//}
//}

//func TestSerializer(t *testing.T) {
//}

//func TestCalculateHash() {
	
//}