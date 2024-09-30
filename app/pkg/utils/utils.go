package utils

import "encoding/base64"

func isBase64Encoded(s string) bool {
	// Base64 string length must be a multiple of 4
	if len(s)%4 != 0 {
		return false
	}
	// Try to decode the string
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}
