package contract

import (
	"encoding/json"
	"fmt"

	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"
	"github.com/liteseed/aogo"
)

const PROCESS = ""

type Contract struct {
	ao     *aogo.AO
	signer *goar.ItemSigner
}

func New(ao *aogo.AO, process string, signer *goar.ItemSigner) *Contract {
	return &Contract{
		ao:     ao,
		signer: signer,
	}
}

func (c *Contract) aoAction(data string, tags []types.Tag) ([]byte, error) {
	mId, err := c.ao.SendMessage(PROCESS, data, tags, "", c.signer)
	if err != nil {
		return nil, err
	}

	result, err := c.ao.ReadResult(PROCESS, mId)
	if err != nil {
		return nil, err
	}
	return result.Messages[0]["Data"].([]byte), nil
}

func (c *Contract) Balance(id string) (string, error) {
	tags := []types.Tag{
		{Name: "Action", Value: "Balance"},
	}
	res, err := c.aoAction(id, tags)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func (c *Contract) Balances() (*map[string]string, error) {
	tags := []types.Tag{
		{Name: "Action", Value: "Balances"},
	}
	res, err := c.aoAction("", tags)
	if err != nil {
		return nil, err
	}
	var balances map[string]string
	err = json.Unmarshal(res, &balances)
	if err != nil {
		return nil, err
	}
	return &balances, nil
}

func (c *Contract) Initiate(dataItemId string, size uint) (*Staker, error) {
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
	tags := []types.Tag{
		{Name: "Action", Value: "Posted"},
		{Name: "Transaction", Value: transactionId},
	}
	_, err := c.ao.SendMessage(PROCESS, "", tags, dataItemId, c.signer)
	return err
}

func (c *Contract) Release(dataItemId string, transactionId string) error {
	tags := []types.Tag{
		{Name: "Action", Value: "Release"},
	}
	_, err := c.ao.SendMessage(PROCESS, "", tags, dataItemId, c.signer)
	return err
}

func (c *Contract) Stake(url string) error {
	_, err := c.ao.SendMessage(PROCESS, "", []types.Tag{{Name: "Action", Value: "Stake"}, {Name: "Url", Value: url}}, "", c.signer)
	return err
}

func (c *Contract) Staked() (string, error) {
	tags := []types.Tag{
		{Name: "Action", Value: "Staked"},
	}
	res, err := c.aoAction("", tags)
	if err != nil {
		return "", err
	}
	return string(res), nil
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
	tags := []types.Tag{
		{Name: "Action", Value: "Transfer"},
		{Name: "Recipient", Value: recipient},
		{Name: "Quantity", Value: quantity},
	}
	_, err := c.ao.SendMessage(PROCESS, "", tags, "", c.signer)
	return err
}

func (c *Contract) Unstake() error {
	_, err := c.ao.SendMessage(PROCESS, "", []types.Tag{{Name: "Action", Value: "Unstake"}}, "", c.signer)
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
