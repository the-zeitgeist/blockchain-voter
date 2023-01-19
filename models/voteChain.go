package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/manifoldco/promptui"
	"github.com/the-zeitgeist/voter/constants"
	"github.com/the-zeitgeist/voter/utils"
)

type VoteChain struct {
	Id           string                 `json:"id"`
	Model        string                 `json:"model"`
	Difficulty   int                    `json:"difficulty"`
	Nodes        []Node                 `json:"nodes"`
	Candidates   []Candidate            `json:"candidates"`
	Chain        []Block                `json:"chain"`
	Txs          map[string]Transaction `json:"txs"`
	IsProcessing bool                   `json:"isProcessing"`
}

func NewVoteChain() (*VoteChain, error) {
	ip, err := utils.GetPublicIp()
	if err != nil {
		fmt.Printf("A public ip address is required to setup the new node\n %s", err.Error())
		return nil, err
	}

	n := Node{Active: true, Ip: ip}
	id := uuid.NewString()
	model := "pluralism"

	var paths []string
	filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".json" {
			paths = append(paths, path)
		}

		return nil
	})

	if len(paths) == 0 {
		return nil, errors.New("candidates json file is required")
	}

	var candidates []Candidate

	prompt := promptui.Select{
		Label: "Select candidates file",
		Items: paths,
	}

	_, path, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Candidates file not found\n")
		return nil, err
	}

	err = json.Unmarshal([]byte(data), &candidates)
	if err != nil {
		fmt.Printf("An error occurred while parsing the candidates file\n")
		return nil, err
	}

	vc := VoteChain{
		Id:           id,
		Model:        model,
		Difficulty:   5,
		Nodes:        []Node{n},
		Candidates:   candidates,
		Chain:        []Block{},
		Txs:          map[string]Transaction{},
		IsProcessing: false,
	}

	fmt.Println("Initiating ...")
	b := NewBlock("Genesis Block", map[string]Transaction{})

	b.Mine(vc.Difficulty)
	vc.Chain = append(vc.Chain, *b)

	err = vc.Export(constants.ConfigFile)

	return &vc, err
}

func JoinVoteChain() (*VoteChain, error) {
	ip, err := utils.GetPublicIp()
	if err != nil {
		fmt.Printf("A public ip address is required to setup the new node\n %s", err.Error())
		return nil, err
	}

	prompt := promptui.Prompt{
		Label: "Vote Chain Address",
		Validate: func(input string) error {
			if len(input) == 0 {
				return errors.New("empty string")
			}

			if !strings.Contains(input, ".") {
				return errors.New("invalid address")
			}

			if input[0] == '.' || input[len(input)-1] == '.' {
				return errors.New("invalid address")
			}

			return nil
		},
	}

	address, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	res, err := http.Get(fmt.Sprintf("%s/join/%s", address, ip))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("couldn't get chain from %s", address)
	}

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &VoteChain{
		IsProcessing: false,
	}, nil
}

func (vc VoteChain) Export(output string) error {
	data, err := json.MarshalIndent(vc, " ", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(output, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (vc *VoteChain) AddTx(tx Transaction) {
	ip, _ := utils.GetPublicIp()
	vc.Txs[tx.Id] = tx

	b, _ := json.Marshal(tx)
	br := bytes.NewReader(b)
	for _, n := range vc.Nodes {
		if n.Ip == ip {
			continue
		}

		go http.Post(fmt.Sprintf("%s/transactions", n.Ip), "application/json", br)
	}

	vc.Export(constants.ConfigFile)
	vc.Proccess()
}

func (vc *VoteChain) AddBlock() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	fmt.Println("Adding...")
	b := NewBlock(vc.Chain[len(vc.Chain)-1].Hash, vc.Txs)

	vc.IsProcessing = true
	b.Transactions = vc.Txs
	vc.Export(constants.ConfigFile)

	b.Mine(vc.Difficulty)
	vc.Chain = append(vc.Chain, *b)
	fmt.Println("Mined", b.Hash)

	ntx := make(map[string]Transaction)
	for k, t := range vc.Txs {
		if _, ok := b.Transactions[k]; !ok {
			ntx[k] = t
		}
	}

	vc.Txs = ntx

	vc.IsProcessing = false
	vc.Export(constants.ConfigFile)
}

func (vc VoteChain) ValidateChain() bool {
	for i, b := range vc.Chain {
		if b.Hash != b.CalculateHash(vc.Difficulty) {
			return false
		}

		if i != len(vc.Chain)-1 {
			if b.Hash != vc.Chain[i+1].PrevHash {
				return false
			}
		}
	}

	return true
}

func (vc *VoteChain) Proccess() {
	if vc.IsProcessing {
		fmt.Println("proccessing, return")
		return
	}

	if len(vc.Txs) == 0 {
		fmt.Println("No txs", vc.Txs)
		return
	}

	go vc.AddBlock()
}

func (vc VoteChain) Result() map[string]int {
	result := make(map[string]int)

	for _, b := range vc.Chain {
		for _, t := range b.Transactions {
			if v, ok := result[t.Candidate]; !ok {
				result[t.Candidate] = 1
			} else {
				result[t.Candidate] = v + 1
			}
		}
	}

	return result
}
