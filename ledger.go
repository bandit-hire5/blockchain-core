package core

import (
	"bufio"
	"encoding/json"
	"os"
	"time"

	"github.com/bandit/blockchain-core/utils"
)

type Ledger struct {
	Filepath string
}

func NewLedger(filepath string) *Ledger {
	return &Ledger{
		Filepath: filepath,
	}
}

func (l *Ledger) CreateWithGenesisBlock() error {
	_, err := os.Stat(l.Filepath)

	if os.IsNotExist(err) {
		file, err := os.Create(l.Filepath)
		if err != nil {
			return err
		}

		defer file.Close()

		block, err := generateGenesisBlock()
		if err != nil {
			return err
		}

		err = l.AddBlock(block)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *Ledger) CreateEmpty() error {
	_, err := os.Stat(l.Filepath)

	if !os.IsNotExist(err) {
		os.Remove(l.Filepath)
	}

	file, err := os.Create(l.Filepath)
	if err != nil {
		return err
	}

	defer file.Close()

	return nil
}

func (l *Ledger) AddBlock(block *Block) error {
	var file, err = os.OpenFile(l.Filepath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer file.Close()

	data, err := json.Marshal(block)
	if err != nil {
		return err
	}

	_, err = file.WriteString(string(data) + "\n")
	if err != nil {
		return err
	}

	err = file.Sync()
	if err != nil {
		return err
	}

	return nil
}

func (l *Ledger) GetAllBlocks() ([]string, error) {
	var list []string

	file, err := os.Open(l.Filepath)
	if err != nil {
		return list, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		list = append(list, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return list, err
	}

	return list, nil
}

func (l *Ledger) GetLastBlock() (*Block, error) {
	var block Block

	list, err := l.GetAllBlocks()
	if err != nil {
		return &block, err
	}

	json.Unmarshal([]byte(list[len(list)-1]), &block)

	return &block, nil
}

func (l *Ledger) NextBlock(data Data) (*Block, error) {
	previousBlock, err := l.GetLastBlock()
	if err != nil {
		return &Block{}, err
	}

	nextIndex := previousBlock.Index + 1
	nextTime := time.Now().UTC()

	dataString, err := json.Marshal(&data)
	if err != nil {
		return &Block{}, err
	}

	nextHash, err := utils.CalculateHash(nextIndex, previousBlock.Hash, nextTime, dataString)
	if err != nil {
		return &Block{}, err
	}

	return &Block{
		Index:        nextIndex,
		Hash:         nextHash,
		PreviousHash: previousBlock.Hash,
		Timestamp:    nextTime.String(),
		Data:         data,
	}, nil
}
