package util

import (
	"bytes"
	"encoding/json"
)

func Keys(words map[string][]string) []string {
	var keys []string

	for key, _ := range words {
		keys = append(keys, key)
	}
	return keys
}

func Perror(err error) {
	if err != nil {
		panic(err)
	}
}

func PrettyPrint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "\t")
	return out.Bytes(), err
}
