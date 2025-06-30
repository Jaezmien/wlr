package jsonl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	A string `json:"a"`
	B string `json:"b"`
}

func TestMarshal(t *testing.T) {
	jsonData := make([]TestStruct, 0)
	jsonData = append(jsonData, TestStruct{A: "a", B: "b"})
	jsonData = append(jsonData, TestStruct{A: "c", B: "d"})

	data, err := Marshal(&jsonData)

	if !assert.NoError(t, err, "Expected no error") {
		return
	}

	assert.Equal(
		t,
		`{"a":"a","b":"b"}`+"\n"+`{"a":"c","b":"d"}`,
		string(data),
		"Expected output not equal",
	)
}

func TestUnmarshal(t *testing.T) {
	data := `{"a":"a","b":"b"}` + "\n" + `{"a":"c","b":"d"}`

	var jsonData []TestStruct
	err := Unmarshal([]byte(data), &jsonData)

	if !assert.NoError(t, err, "Expected no error") {
		return
	}

	assert.Equal(
		t,
		2,
		len(jsonData),
		"Expected two values",
	)
}
