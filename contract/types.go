package contract

type Staker struct {
	ID         string `json:"id"`
	Reputation string `json:"reputation"`
	URL        string `json:"url"`
}

type Upload struct {
	DataItem    string `json:"id"`
	Status      string `json:"status"`
	Bundler     string `json:"bundler"`
	Transaction string `json:"transaction"`
	Paid        string `json:"paid"`
}
