package bittrex

import (
	"encoding/json"
	"time"
)

// Used in getmarkethistory
type Trade struct {
	OrderUuid string
	TimeStamp time.Time
	Quantity  float64
	Price     float64
	Total     float64
	FillType  string
	OrderType string
}

func (t *Trade) UnmarshalJSON(data []byte) (err error) {
	s := struct {
		OrderUuid string  `json:"OrderUuid"`
		TimeStamp string  `json:"TimeStamp"`
		Quantity  float64 `json:"Quantity"`
		Price     float64 `json:"Price"`
		Total     float64 `json:"Total"`
		FillType  string  `json:"FillType"`
		OrderType string  `json:"OrderType"`
	}{}
	if err = json.Unmarshal(data, &s); err != nil {
		return err
	}
	var _t time.Time
	if s.TimeStamp != "" {
		_t, err = time.Parse(TIME_FORMAT, s.TimeStamp)
		if err != nil {
			return err
		}
	}
	*t = Trade{
		OrderUuid: s.OrderUuid,
		TimeStamp: _t,
		Quantity:  s.Quantity,
		Price:     s.Price,
		Total:     s.Total,
		FillType:  s.FillType,
		OrderType: s.OrderType,
	}
	return nil
}

func (t Trade) MarshalJSON() ([]byte, error) {
	_t := ""
	if !t.TimeStamp.IsZero() {
		_t = t.TimeStamp.Format(TIME_FORMAT)
	}
	return json.Marshal(
		struct {
			OrderUuid string  `json:"OrderUuid"`
			TimeStamp string  `json:"TimeStamp"`
			Quantity  float64 `json:"Quantity"`
			Price     float64 `json:"Price"`
			Total     float64 `json:"Total"`
			FillType  string  `json:"FillType"`
			OrderType string  `json:"OrderType"`
		}{
			OrderUuid: t.OrderUuid,
			TimeStamp: _t,
			Quantity:  t.Quantity,
			Price:     t.Price,
			Total:     t.Total,
			FillType:  t.FillType,
			OrderType: t.OrderType,
		})
}
