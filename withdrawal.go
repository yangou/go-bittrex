package bittrex

import (
	"encoding/json"
	"time"
)

type Withdrawal struct {
	PaymentUuid    string
	Currency       string
	Amount         float64
	Address        string
	Opened         time.Time
	Authorized     bool
	PendingPayment bool
	TxCost         float64
	TxId           string
	Canceled       bool
}

func (w *Withdrawal) UnmarshalJSON(data []byte) (err error) {
	s := struct {
		PaymentUuid    string  `json:"PaymentUuid"`
		Currency       string  `json:"Currency"`
		Amount         float64 `json:"Amount"`
		Address        string  `json:"Address"`
		Opened         string  `json:"Opened"`
		Authorized     bool    `json:"Authorized"`
		PendingPayment bool    `json:"PendingPayment"`
		TxCost         float64 `json:"TxCost"`
		TxId           string  `json:"TxId"`
		Canceled       bool    `json:"Canceled"`
	}{}
	if err = json.Unmarshal(data, &s); err != nil {
		return err
	}
	var t time.Time
	if s.Opened != "" {
		t, err = time.Parse(TIME_FORMAT, s.Opened)
		if err != nil {
			return err
		}
	}

	*w = Withdrawal{
		PaymentUuid:    s.PaymentUuid,
		Currency:       s.Currency,
		Amount:         s.Amount,
		Address:        s.Address,
		Opened:         t,
		Authorized:     s.Authorized,
		PendingPayment: s.PendingPayment,
		TxCost:         s.TxCost,
		TxId:           s.TxId,
		Canceled:       s.Canceled,
	}
	return nil
}

func (w Withdrawal) MarshalJSON() ([]byte, error) {
	t := ""
	if !w.Opened.IsZero() {
		t = w.Opened.Format(TIME_FORMAT)
	}
	return json.Marshal(
		struct {
			PaymentUuid    string  `json:"PaymentUuid"`
			Currency       string  `json:"Currency"`
			Amount         float64 `json:"Amount"`
			Address        string  `json:"Address"`
			Opened         string  `json:"Opened"`
			Authorized     bool    `json:"Authorized"`
			PendingPayment bool    `json:"PendingPayment"`
			TxCost         float64 `json:"TxCost"`
			TxId           string  `json:"TxId"`
			Canceled       bool    `json:"Canceled"`
		}{
			PaymentUuid:    w.PaymentUuid,
			Currency:       w.Currency,
			Amount:         w.Amount,
			Address:        w.Address,
			Opened:         t,
			Authorized:     w.Authorized,
			PendingPayment: w.PendingPayment,
			TxCost:         w.TxCost,
			TxId:           w.TxId,
			Canceled:       w.Canceled,
		})
}
