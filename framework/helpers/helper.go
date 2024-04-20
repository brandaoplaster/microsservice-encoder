package helpers

import "encoding/json"

func IsJson(char string) error {
	var js struct{}

	if err := json.Unmarshal([]byte(char), &js); err != nil {
		return err
	}

	return nil
}
