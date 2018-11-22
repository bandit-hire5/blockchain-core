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

func (self *Ledger) CreateWithGenesisBlock() error {
	_, err := os.Stat(self.Filepath)

	if os.IsNotExist(err) {
		file, err := os.Create(self.Filepath)
		if err != nil {
			return err
		}

		defer file.Close()

		block, err := generateGenesisBlock()
		if err != nil {
			return err
		}

		err = self.AddBlock(block)
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *Ledger) CreateEmpty() error {
	_, err := os.Stat(self.Filepath)

	if !os.IsNotExist(err) {
		os.Remove(self.Filepath)
	}

	file, err := os.Create(self.Filepath)
	if err != nil {
		return err
	}

	defer file.Close()

	return nil
}

func (self *Ledger) AddBlock(block *Block) error {
	var file, err = os.OpenFile(self.Filepath, os.O_APPEND|os.O_WRONLY, 0600)
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

func (self *Ledger) GetAllBlocks() ([]string, error) {
	var list []string

	file, err := os.Open(self.Filepath)
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

func (self *Ledger) GetLastBlock() (*Block, error) {
	var block Block

	list, err := self.GetAllBlocks()
	if err != nil {
		return &block, err
	}

	json.Unmarshal([]byte(list[len(list)-1]), &block)

	return &block, nil
}

func (self *Ledger) NextBlock(data Data) (*Block, error) {
	previousBlock, err := self.GetLastBlock()
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
