package utils

import (
	"encoding/json"

	"github.com/samber/lo"
)

func DeepCopyByJson[T any](p T) (T, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return lo.Empty[T](), err
	}
	var newP T
	err = json.Unmarshal(data, &newP)
	if err != nil {
		return lo.Empty[T](), err
	}
	return newP, nil
}

func DeepCopyErrorByJson[T any](p T) T {
	data, err := json.Marshal(p)
	if err != nil {
		panic(err.Error())
	}
	var newP T
	err = json.Unmarshal(data, &newP)
	if err != nil {
		panic(err.Error())
	}
	return newP
}
