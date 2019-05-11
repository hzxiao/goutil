package hex

import "encoding/hex"

func EncodeString(str string) string {
	return hex.EncodeToString([]byte(str))
}

func EncodeBytes(bys []byte) string {
	return hex.EncodeToString(bys)
}

func DecodeHexString(hexStr string) ([]byte, error) {
	return hex.DecodeString(hexStr)
}

func Reverse(hexStr string) string {
	return ""
}


