package blockchain

import (
	"testing"
	"github.com/oliinant/fish-ham/scripts"
)

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

//func TestgetTxListInfoMap() {

//}

//func TestSerializer() {

//}

//func TestCalculateHash() {
	
//}