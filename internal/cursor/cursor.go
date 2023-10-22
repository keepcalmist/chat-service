package cursor

import (
	"encoding/base64"
	"encoding/json"
)

func Encode(data any) (string, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(dataBytes), nil
}

func Decode(in string, to any) error {
	data, err := base64.URLEncoding.DecodeString(in)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, to)
}
