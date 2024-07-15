package main

import (
	"log"

	"github.com/liteseed/goar/signer"
	"github.com/liteseed/goar/tag"
	"github.com/liteseed/goar/transaction/data_item"
	"github.com/liteseed/goar/wallet"
	"github.com/liteseed/sdk-go/contract"
)

const PROCESS = "PWSr59Cf6jxY7aA_cfz69rs0IiJWWbmQA8bAKknHeMo"

func main() {

	bundler, err := signer.FromPath("./examples/bundler.json")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Bundler: %s\n", bundler.Address)

	cb := contract.New(PROCESS, bundler)

	bs, err := cb.Balance(bundler.Address)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(bs)

	info, err := cb.Info()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(info)

	staked, err := cb.Staked()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(staked)

	unstakeResp, err := cb.Unstake()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(unstakeResp)
	staked, err = cb.Staked()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(staked)

	stakeResp, err := cb.Stake("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(stakeResp)

	staked, err = cb.Staked()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(staked)

	user, err := wallet.FromPath("./examples/user.json", "https://arweave.net")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Signer: %s\n", user.Signer.Address)

	c := contract.New(PROCESS, user.Signer)

	balance, err := c.Balance(PROCESS)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Balance: %s\n", balance)

	balances, err := c.Balances()
	if err != nil {
		log.Println(err)
	}
	log.Printf("Balances: %s\n", balances)

	dataItem := data_item.New([]byte{1, 2, 3}, user.Signer.Address, "", &[]tag.Tag{})

	err = dataItem.Sign(user.Signer)
	if err != nil {
		log.Println(err)
	}

	staker, err := c.Initiate(dataItem.ID, len(dataItem.Raw))
	if err != nil {
		log.Println(err)
	}
	log.Println(staker)

	tx := user.CreateTransaction(nil, "", "100000", nil)
	_, err = user.SignTransaction(tx)
	if err != nil {
		log.Println(err)
	}

	err = c.Pay(dataItem.ID, tx.ID)
	if err != nil {
		log.Println(err)
	}

	err = c.Posted(dataItem.ID)
	if err != nil {
		log.Println(err)
	}
}
