package bittrex

type MarketSummary struct {
	MarketName     string  `json:"MarketName"`
	High           float64 `json:"High"`
	Low            float64 `json:"Low"`
	Ask            float64 `json:"Ask"`
	Bid            float64 `json:"Bid"`
	OpenBuyOrders  int     `json:"OpenBuyOrders"`
	OpenSellOrders int     `json:"OpenSellOrders"`
	Volume         float64 `json:"Volume"`
	Last           float64 `json:"Last"`
	BaseVolume     float64 `json:"BaseVolume"`
	PrevDay        float64 `json:"PrevDay"`
	TimeStamp      string  `json:"TimeStamp"`
}

type MarketSummaries []*MarketSummary

func (m MarketSummaries) Len() int           { return len(m) }
func (m MarketSummaries) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m MarketSummaries) Less(i, j int) bool { return m[i].BaseVolume > m[j].BaseVolume }
