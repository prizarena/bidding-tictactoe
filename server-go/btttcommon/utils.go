package btttcommon

import (
	"encoding/base64"
	"encoding/binary"
)

var endian = binary.BigEndian
var base64UrlEncoder = base64.RawURLEncoding

func EncodeID(id int64) string {
	b := make([]byte, 8)
	endian.PutUint64(b, uint64(id))
	return base64UrlEncoder.EncodeToString(b)
}

func DecodeID(s string) (int64, error) {
	b := make([]byte, 8)
	if _, err := base64UrlEncoder.Decode(b, []byte(s)); err != nil {
		return 0, err
	}
	return int64(endian.Uint64(b)), nil
}
