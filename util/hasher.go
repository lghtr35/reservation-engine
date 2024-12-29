package util

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"

	"github.com/lghtr35/reservation-engine/models"
)

type Hasher struct {
	hasher hash.Hash
	secret string
	salt   string
}

func NewHasher(configuration *models.Configuration) (*Hasher, error) {
	h := sha256.New()

	return &Hasher{hasher: h, salt: configuration.GetSalt(), secret: configuration.Secret}, nil
}

func (h *Hasher) GetHash(s string) (string, error) {
	s = fmt.Sprintf("%s%s%s", h.salt, s, h.secret)
	count, err := h.hasher.Write([]byte(s))
	if err != nil || count == 0 {
		return "", err
	}
	res := hex.EncodeToString(h.hasher.Sum(nil))
	h.hasher.Reset()
	return res, nil
}

func (h *Hasher) Verify(hashed string, new string) (bool, error) {
	res, err := h.GetHash(new)
	if err != nil {
		return false, err
	}
	return res == hashed, nil
}
