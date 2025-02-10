package utils

import "encoding/json"

func UnsafeMarshalJSON(v any) string {
	vJSON, _ := json.MarshalIndent(v, "", "  ")
	return string(vJSON)
}
