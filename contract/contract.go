package contract

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"
	"github.com/liteseed/aogo"
)

type Contract struct {
	ao      *aogo.AO
	process string
	signer  *goar.Signer
}

func New(process string, signer *goar.Signer) *Contract {
	ao, err := aogo.New()
	if err != nil {
		log.Fatal(err)
	}

	return &Contract{
		ao:      ao,
		process: process,
		signer:  signer,
	}
}

func (c *Contract) aoAction(data string, tags []types.Tag) ([]byte, error) {
	itemSigner, err := goar.NewItemSigner(c.signer)
	if err != nil {
		return nil, err
	}
	log.Println(c.process)
	mId, err := c.ao.SendMessage(c.process, data, tags, "", itemSigner)
	if err != nil {
		return nil, err
	}

	result, err := c.ao.LoadResult(c.process, mId)
	if err != nil {
		return nil, err
	}
	return []byte(result.Messages[0]["Data"].(string)), nil
}

func (c *Contract) Balance(target string) (string, error) {
	res, err := c.ao.DryRun(aogo.Message{
		Target: c.process,
		Owner:  target,
		Tags:   []types.Tag{{Name: "Action", Value: "Balance"}},
	})
	if err != nil {
		return "", err
	}
	return res.Messages[0]["Data"].(string), nil
}

func (c *Contract) Balances() (*map[string]string, error) {
	tags := []types.Tag{
		{Name: "Action", Value: "Balances"},
	}
	res, err := c.ao.DryRun(aogo.Message{Target: c.process, Owner: c.signer.Address, Tags: tags})
	if err != nil {
		return nil, err
	}
	var balances map[string]string
	err = json.Unmarshal([]byte(res.Messages[0]["Data"].(string)), &balances)
	if err != nil {
		return nil, err
	}
	return &balances, nil
}

func (c *Contract) Initiate(dataItemId string, size int) (*Staker, error) {
	tags := []types.Tag{
		{Name: "Action", Value: "Initiate"},
		{Name: "Size", Value: fmt.Sprint(size)},
	}
	res, err := c.aoAction(dataItemId, tags)
	if err != nil {
		return nil, err
	}
	var staker Staker
	err = json.Unmarshal(res, &staker)
	if err != nil {
		return nil, err
	}
	return &staker, nil
}

func (c *Contract) Posted(dataItemId string, transactionId string) error {
	itemSigner, err := goar.NewItemSigner(c.signer)
	if err != nil {
		return err
	}
	_, err = c.ao.SendMessage(c.process, "", []types.Tag{{Name: "Action", Value: "Posted"}, {Name: "Transaction", Value: transactionId}}, dataItemId, itemSigner)
	return err
}

func (c *Contract) Release(dataItemId string, transactionId string) error {
	itemSigner, err := goar.NewItemSigner(c.signer)
	if err != nil {
		return err
	}
	_, err = c.ao.SendMessage(c.process, "", []types.Tag{{Name: "Action", Value: "Release"}}, dataItemId, itemSigner)
	return err
}

func (c *Contract) Stake(url string) error {
	itemSigner, err := goar.NewItemSigner(c.signer)
	if err != nil {
		return err
	}
	mid, err := c.ao.SendMessage(c.process, "", []types.Tag{{Name: "Action", Value: "Stake"}, {Name: "Url", Value: url}}, "", itemSigner)
	log.Println(mid)
	if err != nil {
		return err
	}
	result, err := c.ao.LoadResult(c.process, mid)
	log.Println(result)
	if err != nil {
		return err
	}
	return err
}

func (c *Contract) Staked() (string, error) {
	tags := []types.Tag{
		{Name: "Action", Value: "Staked"},
	}
	res, err := c.ao.DryRun(aogo.Message{Target: c.process, Data: c.signer.Address, Owner: c.signer.Address, Tags: tags})
	if err != nil {
		return "", err
	}
	return res.Messages[0]["Data"].(string), nil
}

func (c *Contract) Stakers() (*[]Staker, error) {
	tags := []types.Tag{
		{Name: "Action", Value: "Stakers"},
	}
	res, err := c.aoAction("", tags)
	if err != nil {
		return nil, err
	}
	var stakers []Staker
	err = json.Unmarshal(res, &stakers)
	if err != nil {
		return nil, err
	}
	return &stakers, nil
}

func (c *Contract) Transfer(recipient string, quantity string) error {
	itemSigner, err := goar.NewItemSigner(c.signer)
	if err != nil {
		return err
	}
	_, err = c.ao.SendMessage(c.process, "", []types.Tag{{Name: "Action", Value: "Transfer"}, {Name: "Recipient", Value: recipient}, {Name: "Quantity", Value: quantity}}, "", itemSigner)
	return err
}

func (c *Contract) Unstake() error {
	itemSigner, err := goar.NewItemSigner(c.signer)
	if err != nil {
		return err
	}
	_, err = c.ao.SendMessage(c.process, "", []types.Tag{{Name: "Action", Value: "Unstake"}}, "", itemSigner)
	return err
}

func (c *Contract) Upload(id string) (*Upload, error) {
	tags := []types.Tag{
		{Name: "Action", Value: "Stakers"},
	}
	res, err := c.aoAction(id, tags)
	if err != nil {
		return nil, err
	}
	var upload Upload
	err = json.Unmarshal(res, &upload)
	if err != nil {
		return nil, err
	}
	return &upload, nil
}
