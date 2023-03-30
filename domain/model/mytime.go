package model

import (
	"encoding/json"
	"strings"
	"time"
)

type MyTime struct {
	*time.Time
}

// UnmarshalJSON will be called by json.Unmarshal
func (mt *MyTime) UnmarshalJSON(data []byte) error {
	// 要素がそのまま渡ってくるので "(ダブルクォート)でQuoteされてる
	s := strings.Trim(string(data), "\"")
	if s == "null" {
		t, err := time.Parse("2006/01/02", string(data))
		*mt = MyTime{&t}
		return err
	}
	t, err := time.Parse("2006/01/02", s)
	*mt = MyTime{&t}
	return err
}

// MarshalJSON will be called by json.Marshal
func (mt MyTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(mt.Format("2006/01/02"))
}
