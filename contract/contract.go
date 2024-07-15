package contract

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/liteseed/aogo"
	"github.com/liteseed/goar/signer"
	"github.com/liteseed/goar/tag"
)

type IContract interface {
	Info() (*Info, error)
	Balance(target string) (string, error)
	Balances() (*map[string]string, error)
	Transfer(recipient string, quantity string) error
	Initiate(dataItemID string, size int) (*Staker, error)
	Pay(ID string, paymentID string) error
	Posted(dataItemID string) error
	Release(dataItemID string) error
	Stake(url string) (string, error)
	Unstake() (string, error)
	Staked() (string, error)
	Stakers() (*[]Staker, error)
}

type Contract struct {
	ao      *aogo.AO
	process string
	signer  *signer.Signer
}

func New(process string, signer *signer.Signer) *Contract {
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

func Custom(ao *aogo.AO, process string, signer *signer.Signer) *Contract {
	return &Contract{
		ao:      ao,
		process: process,
		signer:  signer,
	}
}

func (c *Contract) aoAction(data string, tags *[]tag.Tag) ([]byte, error) {
	mID, err := c.ao.SendMessage(c.process, data, tags, "", c.signer)
	if err != nil {
		return nil, err
	}

	result, err := c.ao.LoadResult(c.process, mID)
	if err != nil {
		return nil, err
	}
	return []byte(result.Messages[0]["Data"].(string)), nil
}

func (c *Contract) Info() (*Info, error) {
	res, err := c.ao.DryRun(
		aogo.Message{
			Target: c.process,
			Owner:  c.signer.Address,
			Tags:   &[]tag.Tag{{Name: "Action", Value: "Info"}},
		},
	)
	if err != nil {
		return nil, err
	}
	var info Info
	err = json.Unmarshal([]byte(res.Messages[0]["Data"].(string)), &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (c *Contract) Balance(target string) (string, error) {
	res, err := c.ao.DryRun(aogo.Message{
		Target: c.process,
		Owner:  target,
		Tags:   &[]tag.Tag{{Name: "Action", Value: "Balance"}},
	})
	if err != nil {
		return "", err
	}
	b := res.Messages[0]["Data"]
	if b == nil {
		return "0", nil
	} else {
		return res.Messages[0]["Data"].(string), nil
	}
}

func (c *Contract) Balances() (*map[string]string, error) {
	tags := &[]tag.Tag{
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

func (c *Contract) Initiate(dataItemID string, size int) (*Staker, error) {
	tags := &[]tag.Tag{
		{Name: "Action", Value: "Initiate"},
		{Name: "Size", Value: fmt.Sprint(size)},
	}
	res, err := c.aoAction(dataItemID, tags)
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

func (c *Contract) Pay(ID string, paymentID string) error {
	_, err := c.ao.SendMessage(c.process, ID, &[]tag.Tag{{Name: "Action", Value: "Pay"}, {Name: "Payment", Value: paymentID}}, "", c.signer)
	return err
}

func (c *Contract) Posted(dataItemID string) error {
	_, err := c.ao.SendMessage(c.process, dataItemID, &[]tag.Tag{{Name: "Action", Value: "Posted"}}, "", c.signer)
	return err
}

func (c *Contract) Release(dataItemID string) error {
	_, err := c.ao.SendMessage(c.process, dataItemID, &[]tag.Tag{{Name: "Action", Value: "Release"}}, dataItemID, c.signer)
	return err
}

func (c *Contract) Stake(url string) (string, error) {
	mID, err := c.ao.SendMessage(c.process, "", &[]tag.Tag{{Name: "Action", Value: "Stake"}, {Name: "Url", Value: url}}, "", c.signer)
	if err != nil {
		return "", err
	}
	result, err := c.ao.LoadResult(c.process, mID)
	if err != nil {
		return "", err
	}
	if result.Error != "" {
		return "", fmt.Errorf(result.Error)
	}
	return result.Messages[0]["Data"].(string), nil
}

func (c *Contract) Staked() (string, error) {
	tags := &[]tag.Tag{
		{Name: "Action", Value: "Staked"},
	}
	res, err := c.ao.DryRun(aogo.Message{Target: c.process, Data: c.signer.Address, Owner: c.signer.Address, Tags: tags})
	if err != nil {
		return "", err
	}
	return res.Messages[0]["Data"].(string), nil
}

func (c *Contract) Stakers() (*[]Staker, error) {
	tags := &[]tag.Tag{
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
	_, err := c.ao.SendMessage(c.process, "", &[]tag.Tag{{Name: "Action", Value: "Transfer"}, {Name: "Recipient", Value: recipient}, {Name: "Quantity", Value: quantity}}, "", c.signer)
	return err
}

func (c *Contract) Unstake() (string, error) {
	mID, err := c.ao.SendMessage(c.process, "", &[]tag.Tag{{Name: "Action", Value: "Unstake"}}, "", c.signer)
	if err != nil {
		return "", err
	}

	result, err := c.ao.LoadResult(c.process, mID)
	if err != nil {
		return "", err
	}
	if result.Error != "" {
		return "", fmt.Errorf(result.Error)
	}
	return result.Messages[0]["Data"].(string), nil
}
