package jsonl

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
)

func Marshal[T any](v *[]T) ([]byte, error) {
	data := make([][]byte, 0)

	for _, x := range *v {
		e, err := json.Marshal(x)
		if err != nil {
			return []byte{}, fmt.Errorf("error while trying to unmarshal line: %v", err)
		}
		data = append(data, e)
	}

	return bytes.Join(data, []byte("\n")), nil
}

func Unmarshal[T any](data []byte, v *[]T) error {
	scanner := bufio.NewScanner(bytes.NewReader(data))

	for scanner.Scan() {
		line := scanner.Text()

		var entry T
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			return fmt.Errorf("error while trying to unmarshal line: %v", err)
		}

		*v = append(*v, entry)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error while trying to unmarshal data: %v", err)
	}

	return nil
}
