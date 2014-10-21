package util

import (
	"errors"
)

func ByteToBase10(b []byte) (n uint, err error) {
	base := uint(10)

	n = 0
	for i := 0; i < len(b); i++ {
		var v byte
		d := b[i]
		switch {
		case '0' <= d && d <= '9':
			v = d - '0'
		default:
			n = 0
			err = errors.New("failed to convert to Base10")
			break
		}
		n *= base
		n += uint(v)
	}

	return n, err
}
