package bittrex

import (
	"encoding/json"
	"time"
)

type Deposit struct {
	Id            int64
	Amount        float64
	Currency      string
	Confirmations int
	LastUpdated   time.Time
	TxId          string
	CryptoAddress string
}

func (d *Deposit) UnmarshalJSON(data []byte) (err error) {
	s := struct {
		Id            int64   `json:"Id"`
		Amount        float64 `json:"Amount"`
		Currency      string  `json:"Currency"`
		Confirmations int     `json:"Confirmations"`
		LastUpdated   string  `json:"LastUpdated"`
		TxId          string  `json:"TxId"`
		CryptoAddress string  `json:"CryptoAddress"`
	}{}
	if err = json.Unmarshal(data, &s); err != nil {
		return err
	}
	var t time.Time
	if s.LastUpdated != "" {
		t, err = time.Parse(TIME_FORMAT, s.LastUpdated)
		if err != nil {
			return err
		}
	}

	*d = Deposit{
		Id:            s.Id,
		Amount:        s.Amount,
		Currency:      s.Currency,
		Confirmations: s.Confirmations,
		LastUpdated:   t,
		TxId:          s.TxId,
		CryptoAddress: s.CryptoAddress,
	}
	return nil
}

func (d Deposit) MarshalJSON() ([]byte, error) {
	t := ""
	if !d.LastUpdated.IsZero() {
		t = d.LastUpdated.Format(TIME_FORMAT)
	}

	return json.Marshal(
		struct {
			Id            int64   `json:"Id"`
			Amount        float64 `json:"Amount"`
			Currency      string  `json:"Currency"`
			Confirmations int     `json:"Confirmations"`
			LastUpdated   string  `json:"LastUpdated"`
			TxId          string  `json:"TxId"`
			CryptoAddress string  `json:"CryptoAddress"`
		}{
			Id:            d.Id,
			Amount:        d.Amount,
			Currency:      d.Currency,
			Confirmations: d.Confirmations,
			LastUpdated:   t,
			TxId:          d.TxId,
			CryptoAddress: d.CryptoAddress,
		})
}
