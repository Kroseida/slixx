package dto

import (
	"encoding/json"
)

func Map(from interface{}, to interface{}) {
	data, err := json.Marshal(from)

	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, to)
	if err != nil {
		panic(err)
	}
}
