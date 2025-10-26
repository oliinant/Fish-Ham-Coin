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

func InitCustomTx(time time.Time, sender string, receiver string, amount float64, taxPercent float64) *Transaction {
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
func InitDefaultTx(sender string, receiver string, amount float64) *Transaction {
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
	Limit int
}

func (txs *TxList) TxListInfoMap() []map[string]interface{} {
	txListMap := []map[string]interface{}{}

	for _, tx := range txs.Transactions {
		info := scripts.InfoMap(tx)
		txListMap = append(txListMap, info)
	}
	return txListMap
}

func (txs *TxList) LengthLimiter() (bool, error) {
	if txs.Limit <= len(txs.Transactions) {
		return true, fmt.Errorf("Maxium number of Transactions reached")
	}
	return false, nil
}

func (txs *TxList) AddTransaction(tx *Transaction) error {
	if txs.Transactions == nil {
		txs.Transactions = make(map[uuid.UUID]*Transaction)
	}
	
	state, err := txs.LengthLimiter()
	if state {
		return scripts.WrapError("Failed to add transaction", err)
	}

	txs.Transactions[tx.ID] = tx
	return nil
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
	err := scripts.checkTextExistence(s)
	if err != nil {
		return nil, scripts.WrapError("Failed to convert hash", err)
	}

	for _, r  := range s {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F' )) {
			return "", fmt.Errorf("Invalid character in hash: %q", r)
		}
	}
	return Hash(s), nil
}

type PoW struct {
	Nonce int `json:"nonce"`
	Difficulty int `json:"difficulty"`
	Hash Hash `json:"hash"`
}

type Block struct {
	Index int `json:"index"`
	Timestamp time.Time `json:"timestamp"`
	Transactions TxList `json:"transactions"`
	PrevHash Hash `json:"prev_hash"`
	PoW: PoW `json:"proof_of_work"`
	Reward float64 `json:"reward"`
	TotalDifficulty int `json:"total_difficulty"`
}

func NewBlock(prevHash Hash, txList TxList, nonce int) (*Block, error) {
	prevBlock, err := BlockByHash("Chain", prevHash)
	if err != nil {
		return nil, scripts.WrapError("Block creation failed", err)
	}

	block := &Block{
		Index: prevBlock.Index + 1,
		Timestamp: time.Now(),
		Transactions: txList,
		PrevHash: prevHash,
		PoW: PoW{
			Nonce: nonce
			Hash: "",
		Reward: 0,
		TotalDifficulty: prevBlock.TotalDifficulty,
		}
	}
}

func (b *Block) ProccessBlock() erorr{
	b.PoW.Hash = b.CalculateHash()

	err = scripts.checkTextExistence(b.PoW.Hash)
	if err != nil {
		return scripts.WrapError("Failed to process block")
	}
	b.PoW.Difficulty = b.CalculateDifficulty()
}

func (b *Block) Reward() float64 {

}

func (b *Block) Serializer() (string, error) {
	delimiter := "\\fh~8"

	dataJSON, err := json.Marshal(b.Transactions)
	if err != nil {
		return "", scripts.WrapError("Failed to serialize block", err)
	}

	serializedBlock := fmt.Sprintf(
		"%d%s%s%s%s%s%s%d",
		b.Index, delimiter,
		b.Timestamp, delimiter,
		string(dataJSON), delimiter,
		b.PrevHash, delimiter,
		b.PoW.Nonce,
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

func (b *Block)CalculateDifficulty() (int, error) {
	err = scripts.checkTextExistence()
	if err != nil {
		return scripts.WrapError("Failed to calculate difficulty", err)
	}

	difficulty := 0
	hexZeroConversions := map[string]int{
		"0": 4,
		"1": 3,
		"2": 2,
		"3": 2,
		"4": 1,
		"5": 1,
		"6": 1,
		"7": 1,
	}

	for _,r := range b.PoW.Hash {
		s := string(r)
		value, ok := hexZeroConversions[s]

		if !ok {
			return difficulty
		} else if s == "0" {
			difficulty += value
		} else {
			return difficulty + value
		}
	}
	return difficulty, nil
}


type Blockchain struct {
	Chain map[Hash]*Block `json:"blockchain"`
	Tips map[Hash]*Block `json:"canonical_branch"`
	Genesis *Block `json:genesis_block`
}

func (bc *Blockchain) InitGenesis() {
	bc.Genesis = &Block{
		Index: 0,
		Timestamp: time.Now(),
		Transactions: make(map[uuid.UUID]*Transaction),
		PrevHash: Hash("")
		PoW: PoW{
			Nonce: 0,
			Difficulty: 0,
			ash: Hash("a0324dc6955ab9ec2ffb7bc1922b81d5cd0bba429ee427fe1bcba94ee74d7a13"),
		}
		Reward: 0,
		ChainDifficulty: 0
	}

	bc.Chain = map[Hash]*Block{bc.Genesis.Hash: bc.Genesis}
	bc.Tips = map[Hash]*Block{bc.Genesis.Hash: bc.Genesis}
}

func (bc *Blockchain) BlockByHash(mapType string, hash Hash) (*Block, error) {
	var block *Block
	var ok bool

	if mapType == "Chain" {
		block, ok := bc.Chain[hash]
	} else if mapType == "Tips" {
		block, ok := bc.Tips[hash]
	} else {
		return nil, fmt.Errorf("Invalid map %s specified", mapType)
	}

	if !ok {
		return nil, fmt.Errorf("Block %s not found", hash)
	}
	return block, nil
}

func (bc *Blockchain) AddBlock() {
	
}