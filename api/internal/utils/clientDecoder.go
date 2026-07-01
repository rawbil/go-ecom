package utils

import (
	"encoding/json"
	"net/http"
)

func DecodeClient(r *http.Request, body any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(body)
}
