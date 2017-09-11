package bittrex

import (
	"encoding/json"
	"time"
)

type Interval string

const (
	OneMin    Interval = "oneMin"
	FiveMin   Interval = "fiveMin"
	ThirtyMin Interval = "thirtyMin"
	Hour      Interval = "hour"
	Day       Interval = "day"
)

var CANDLE_INTERVALS = map[Interval]bool{
	OneMin:    true,
	FiveMin:   true,
	ThirtyMin: true,
	Hour:      true,
	Day:       true,
}

type Candle struct {
	TimeStamp  time.Time
	Open       float64
	Close      float64
	High       float64
	Low        float64
	Volume     float64
	BaseVolume float64
}

type NewCandles struct {
	Ticks []Candle `json:"ticks"`
}

func (c *Candle) UnmarshalJSON(data []byte) (err error) {
	s := struct {
		TimeStamp  string  `json:"T"`
		Open       float64 `json:"O"`
		Close      float64 `json:"C"`
		High       float64 `json:"H"`
		Low        float64 `json:"L"`
		Volume     float64 `json:"V"`
		BaseVolume float64 `json:"BV"`
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

	*c = Candle{
		TimeStamp:  t,
		Open:       s.Open,
		Close:      s.Close,
		High:       s.High,
		Low:        s.Low,
		Volume:     s.Volume,
		BaseVolume: s.BaseVolume,
	}
	return nil
}

func (c Candle) MarshalJSON() ([]byte, error) {
	t := ""
	if !c.TimeStamp.IsZero() {
		t = c.TimeStamp.Format(TIME_FORMAT)
	}
	return json.Marshal(
		struct {
			TimeStamp  string  `json:"T"`
			Open       float64 `json:"O"`
			Close      float64 `json:"C"`
			High       float64 `json:"H"`
			Low        float64 `json:"L"`
			Volume     float64 `json:"V"`
			BaseVolume float64 `json:"BV"`
		}{
			TimeStamp:  t,
			Open:       c.Open,
			Close:      c.Close,
			High:       c.High,
			Low:        c.Low,
			Volume:     c.Volume,
			BaseVolume: c.BaseVolume,
		})
}
