package bittrex

import (
	"encoding/json"
	"time"
)

type Order struct {
	OrderUuid         string
	Exchange          string
	TimeStamp         time.Time
	OrderType         string
	Limit             float64
	Quantity          float64
	QuantityRemaining float64
	Commission        float64
	Price             float64
	PricePerUnit      float64
}

// For getorder
type Order2 struct {
	AccountId                  string
	OrderUuid                  string `json:"OrderUuid"`
	Exchange                   string `json:"Exchange"`
	Type                       string
	Quantity                   float64 `json:"Quantity"`
	QuantityRemaining          float64 `json:"QuantityRemaining"`
	Limit                      float64 `json:"Limit"`
	Reserved                   float64
	ReserveRemaining           float64
	CommissionReserved         float64
	CommissionReserveRemaining float64
	CommissionPaid             float64
	Price                      float64 `json:"Price"`
	PricePerUnit               float64 `json:"PricePerUnit"`
	Opened                     string
	Closed                     string
	IsOpen                     bool
	Sentinel                   string
	CancelInitiated            bool
	ImmediateOrCancel          bool
	IsConditional              bool
	Condition                  string
	ConditionTarget            string
}

func (o *Order) UnmarshalJSON(data []byte) (err error) {
	s := struct {
		OrderUuid         string  `json:"OrderUuid"`
		Exchange          string  `json:"Exchange"`
		TimeStamp         string  `json:"TimeStamp"`
		OrderType         string  `json:"OrderType"`
		Limit             float64 `json:"Limit"`
		Quantity          float64 `json:"Quantity"`
		QuantityRemaining float64 `json:"QuantityRemaining"`
		Commission        float64 `json:"Commission"`
		Price             float64 `json:"Price"`
		PricePerUnit      float64 `json:"PricePerUnit"`
	}{}
	if err = json.Unmarshal(data, &s); err != nil {
		return err
	}
	var t time.Time
	if s.TimeStamp != "" {
		t, err = time.Parse(TIME_FORMAT, s.TimeStamp)
		if err != nil {
			return err
		}
	}

	*o = Order{
		OrderUuid:         s.OrderUuid,
		Exchange:          s.Exchange,
		TimeStamp:         t,
		OrderType:         s.OrderType,
		Limit:             s.Limit,
		Quantity:          s.Quantity,
		QuantityRemaining: s.QuantityRemaining,
		Commission:        s.Commission,
		Price:             s.Price,
		PricePerUnit:      s.PricePerUnit,
	}
	return nil
}

func (o Order) MarshalJSON() ([]byte, error) {
	t := ""
	if !o.TimeStamp.IsZero() {
		t = o.TimeStamp.Format(TIME_FORMAT)
	}
	return json.Marshal(
		struct {
			OrderUuid         string  `json:"OrderUuid"`
			Exchange          string  `json:"Exchange"`
			TimeStamp         string  `json:"TimeStamp"`
			OrderType         string  `json:"OrderType"`
			Limit             float64 `json:"Limit"`
			Quantity          float64 `json:"Quantity"`
			QuantityRemaining float64 `json:"QuantityRemaining"`
			Commission        float64 `json:"Commission"`
			Price             float64 `json:"Price"`
			PricePerUnit      float64 `json:"PricePerUnit"`
		}{
			OrderUuid:         o.OrderUuid,
			Exchange:          o.Exchange,
			TimeStamp:         t,
			OrderType:         o.OrderType,
			Limit:             o.Limit,
			Quantity:          o.Quantity,
			QuantityRemaining: o.QuantityRemaining,
			Commission:        o.Commission,
			Price:             o.Price,
			PricePerUnit:      o.PricePerUnit,
		})
}
