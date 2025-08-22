package jsonstring

import (
	"encoding/json"
	"fmt"
)

// JSONString is a custom type for handling JSON strings
type JSONString string

// UnmarshalJSON implements the json.Unmarshaler interface for JSONString
func (js *JSONString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}
	var temp map[string]interface{}
	err := json.Unmarshal(data, &temp)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSONString: %s", err)
	}

	normalizedData, err := json.Marshal(temp)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	*js = JSONString(normalizedData)
	return nil
}

func (js JSONString) MarshalJSON() ([]byte, error) {
	var temp map[string]interface{}
	err := json.Unmarshal([]byte(js), &temp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSONString: %s", err)
	}

	return json.Marshal(temp)
}
