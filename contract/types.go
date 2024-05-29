package contract

type Staker struct {
	ID         string `json:"id"`
	Reputation int    `json:"reputation"`
	URL        string `json:"url"`
}

type Upload struct {
	DataItem    string `json:"id"`
	Status      string `json:"status"`
	Bundler     string `json:"bundler"`
	Transaction string `json:"transaction"`
	Paid        string `json:"paid"`
}

type Info struct {
	Target       string `json:"Target"`
	Name         string `json:"Name"`
	Ticker       string `json:"Ticker"`
	Logo         string `json:"Logo"`
	Denomination string `json:"Denomination"`
}
