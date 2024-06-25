package contract

import (
	"encoding/json"
	"testing"

	"github.com/liteseed/aogo"
	"github.com/liteseed/goar/signer"
	"github.com/liteseed/goar/tag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock structures:
type MockAO struct {
	mock.Mock
	aogo.AO
}

func (m *MockAO) SendMessage(process, data string, tags *[]tag.Tag, paymentID string, signer *signer.Signer) (string, error) {
	args := m.Called(process, data, tags, paymentID, signer)
	return args.String(0), args.Error(1)
}

func (m *MockAO) LoadResult(process, messageID string) (*aogo.Response, error) {
	args := m.Called(process, messageID)
	return args.Get(0).(*aogo.Response), args.Error(1)
}

func (m *MockAO) DryRun(message aogo.Message) (*aogo.Response, error) {
	args := m.Called(message)
	return args.Get(0).(*aogo.Response), args.Error(1)
}

// Tests
func TestInfo(t *testing.T) {
	mockAO := new(MockAO)
	mockSigner := new(signer.Signer)

	contract := &Contract{
		ao:      &mockAO.AO,
		process: "test_process",
		signer:  mockSigner,
	}

	expectedInfo := Info{
		Name: "Test Contract",
	}

	mockResult := &aogo.Response{
		Messages: []map[string]interface{}{
			{"Data": string(mustMarshal(expectedInfo))},
		},
	}

	mockAO.On("DryRun", mock.Anything).Return(mockResult, nil)

	info, err := contract.Info()
	assert.NoError(t, err)
	assert.Equal(t, &expectedInfo, info)
}

func TestBalance(t *testing.T) {
	mockAO := new(MockAO)
	mockSigner := new(signer.Signer)

	contract := &Contract{
		ao:      &mockAO.AO,
		process: "test_process",
		signer:  mockSigner,
	}

	mockResult := &aogo.Response{
		Messages: []map[string]interface{}{
			{"Data": "1000"},
		},
	}

	mockAO.On("DryRun", mock.Anything).Return(mockResult, nil)

	balance, err := contract.Balance("test_target")
	assert.NoError(t, err)
	assert.Equal(t, "1000", balance)
}

func TestBalances(t *testing.T) {
	mockAO := new(MockAO)
	mockSigner := new(signer.Signer)

	contract := &Contract{
		ao:      &mockAO.AO,
		process: "test_process",
		signer:  mockSigner,
	}

	expectedBalances := map[string]string{
		"address1": "1000",
		"address2": "500",
	}

	mockResult := &aogo.Response{
		Messages: []map[string]interface{}{
			{"Data": string(mustMarshal(expectedBalances))},
		},
	}

	mockAO.On("DryRun", mock.Anything).Return(mockResult, nil)

	balances, err := contract.Balances()
	assert.NoError(t, err)
	assert.Equal(t, &expectedBalances, balances)
}

func TestInitiate(t *testing.T) {
	mockAO := new(MockAO)
	mockSigner := new(signer.Signer)

	contract := &Contract{
		ao:      &mockAO.AO,
		process: "test_process",
		signer:  mockSigner,
	}

	expectedStaker := Staker{
		ID: "staker1",
	}

	mockAO.On("SendMessage", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("messageID", nil)
	mockAO.On("LoadResult", mock.Anything, mock.Anything).Return(&aogo.Response{
		Messages: []map[string]interface{}{
			{"Data": string(mustMarshal(expectedStaker))},
		},
	}, nil)

	staker, err := contract.Initiate("dataItemID", 100)
	assert.NoError(t, err)
	assert.Equal(t, &expectedStaker, staker)
}

func TestPay(t *testing.T) {
	mockAO := new(MockAO)
	mockSigner := new(signer.Signer)

	contract := &Contract{
		ao:      &mockAO.AO,
		process: "test_process",
		signer:  mockSigner,
	}

	mockAO.On("SendMessage", "test_process", "dataItemID", &[]tag.Tag{{Name: "Action", Value: "Pay"}, {Name: "Payment", Value: "paymentID"}}, "", mockSigner).Return("messageID", nil)

	err := contract.Pay("dataItemID", "paymentID")
	assert.NoError(t, err)
}

func TestPosted(t *testing.T) {
	mockAO := new(MockAO)
	mockSigner := new(signer.Signer)

	contract := &Contract{
		ao:      &mockAO.AO,
		process: "test_process",
		signer:  mockSigner,
	}

	mockAO.On("SendMessage", "test_process", "dataItemID", &[]tag.Tag{{Name: "Action", Value: "Posted"}}, "", mockSigner).Return("messageID", nil)

	err := contract.Posted("dataItemID")
	assert.NoError(t, err)
}

func TestRelease(t *testing.T) {
	mockAO := new(MockAO)
	mockSigner := new(signer.Signer)

	contract := &Contract{
		ao:      &mockAO.AO,
		process: "test_process",
		signer:  mockSigner,
	}

	mockAO.On("SendMessage", "test_process", "dataItemID", &[]tag.Tag{{Name: "Action", Value: "Release"}}, "dataItemID", mockSigner).Return("messageID", nil)

	err := contract.Release("dataItemID")
	assert.NoError(t, err)
}

func TestStake(t *testing.T) {
	mockAO := new(MockAO)
	mockSigner := new(signer.Signer)

	contract := &Contract{
		ao:      &mockAO.AO,
		process: "test_process",
		signer:  mockSigner,
	}

	mockAO.On("SendMessage", "test_process", "", &[]tag.Tag{{Name: "Action", Value: "Stake"}, {Name: "Url", Value: "http://example.com"}}, "", mockSigner).Return("messageID", nil)
	mockAO.On("LoadResult", "test_process", "messageID").Return(&aogo.Response{
		Messages: []map[string]interface{}{
			{"Data": "staked"},
		},
	}, nil)

	staked, err := contract.Stake("http://example.com")
	assert.NoError(t, err)
	assert.Equal(t, "staked", staked)
}

func TestUnstake(t *testing.T) {
	mockAO := new(MockAO)
	mockSigner := new(signer.Signer)

	contract := &Contract{
		ao:      &mockAO.AO,
		process: "test_process",
		signer:  mockSigner,
	}

	mockAO.On("SendMessage", "test_process", "", &[]tag.Tag{{Name: "Action", Value: "Unstake"}}, "", mockSigner).Return("messageID", nil)
	mockAO.On("LoadResult", "test_process", "messageID").Return(&aogo.Response{
		Messages: []map[string]interface{}{
			{"Data": "unstaked"},
		},
	}, nil)

	unstaked, err := contract.Unstake()
	assert.NoError(t, err)
	assert.Equal(t, "unstaked", unstaked)
}

func TestStaked(t *testing.T) {
	mockAO := new(MockAO)
	mockSigner := new(signer.Signer)

	contract := &Contract{
		ao:      &mockAO.AO,
		process: "test_process",
		signer:  mockSigner,
	}

	mockResult := &aogo.Response{
		Messages: []map[string]interface{}{
			{"Data": "staked"},
		},
	}

	mockAO.On("DryRun", mock.Anything).Return(mockResult, nil)

	staked, err := contract.Staked()
	assert.NoError(t, err)
	assert.Equal(t, "staked", staked)
}

func TestStakers(t *testing.T) {
	mockAO := new(MockAO)
	mockSigner := new(signer.Signer)

	contract := &Contract{
		ao:      &mockAO.AO,
		process: "test_process",
		signer:  mockSigner,
	}

	expectedStakers := []Staker{
		{ID: "staker1"},
		{ID: "staker2"},
	}

	mockAO.On("SendMessage", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("messageID", nil)
	mockAO.On("LoadResult", mock.Anything, mock.Anything).Return(&aogo.Response{
		Messages: []map[string]interface{}{
			{"Data": string(mustMarshal(expectedStakers))},
		},
	}, nil)

	stakers, err := contract.Stakers()
	assert.NoError(t, err)
	assert.Equal(t, &expectedStakers, stakers)
}

// helper func
func mustMarshal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}
