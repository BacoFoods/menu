package shared

import "encoding/base64"

func Base64Decode(b64 string) ([]byte, error) {
	decoded := base64.StdEncoding.DecodedLen(len(b64))
	decodedString := make([]byte, decoded)
	n, err := base64.StdEncoding.Decode(decodedString, []byte(b64))
	if err != nil {
		return nil, err
	}
	return decodedString[:n], nil
}
