package goutil

import (
	"crypto/md5"
	"fmt"
)

func MD5(bys []byte) string {
	h := md5.New()
	h.Write(bys)
	return fmt.Sprintf("%x", h.Sum(nil))
}

