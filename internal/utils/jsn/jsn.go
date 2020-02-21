package jsn

import "encoding/json"

func B(inp interface{}) string {
	b, err := json.MarshalIndent(&inp, "", "   ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}
