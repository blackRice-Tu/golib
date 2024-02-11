package xcommon

import (
	"encoding/base64"
	"encoding/json"

	"github.com/blackRice-Tu/golib"
)

func JsonMarshal(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		logger := golib.GetStdLogger()
		logger.Println(err)
		return ""
	}
	return string(b)
}

func JsonConvert(from any, to any) error {
	if from == nil {
		return nil
	}
	b, err := json.Marshal(from)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, to)
}

func JsonMarshalToBase64(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func JsonUnmarshalFromBase64(s string, body any) error {
	dataBytes, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return err
	}
	return json.Unmarshal(dataBytes, body)
}
