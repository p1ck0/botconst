package hash

import (
	"crypto/md5"
	"fmt"
)

type Hasher interface {
	Hash(data string) (string, error)
}

type MD5Hasher struct {
	salt string
}

func NewMD5Hasher(salt string) *MD5Hasher {
	return &MD5Hasher{salt: salt}
}

func (h *MD5Hasher) Hash(data string) (string, error) {
	hasher := md5.New()

	if _, err := hasher.Write([]byte(data)); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hasher.Sum([]byte(h.salt))), nil
}
