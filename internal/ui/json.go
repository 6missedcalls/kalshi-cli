package ui

import (
	"encoding/json"
	"fmt"
	"os"
)

func PrintJSON(v interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(v)
}

func PrintJSONCompact(v interface{}) error {
	return json.NewEncoder(os.Stdout).Encode(v)
}

func ToJSONString(v interface{}) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func PrintPlain(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

func Output(format OutputFormat, tableFunc func(), jsonData interface{}, plainFunc func()) error {
	switch format {
	case FormatJSON:
		return PrintJSON(jsonData)
	case FormatPlain:
		plainFunc()
		return nil
	default:
		tableFunc()
		return nil
	}
}
