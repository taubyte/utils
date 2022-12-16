package network

import (
	"encoding/binary"
	"errors"
)

/* Converts uint64 to []byte */
func UInt64ToBytes(i uint64) []byte {
	bs := make([]byte, 8)
	// TCP-IP used BigEndian
	binary.BigEndian.PutUint64(bs, i)
	return bs
}

/* Converts []byte to uint64 */
func BytesToUInt64(data []byte) (uint64, error) {
	if data == nil || len(data) != 8 {
		return 0, errors.New("Invalid data")
	}
	// TCP-IP used BigEndian
	return binary.BigEndian.Uint64(data), nil
}
