package util

import (
	"math/rand/v2"
)

func GetRandString(length int) string {
	alphanumericals := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	maxRand := len(alphanumericals)
	res := make([]rune, length)
	for i := 0; i < length; i++ {
		res[i] = alphanumericals[rand.IntN(maxRand)]
	}

	return string(res)
}

func GetRandHexString(elemCount int) string {
	hexSymbols := []rune("abcdef0123456789")
	maxRand := len(hexSymbols)
	res := make([]rune, elemCount)
	for i := 0; i < elemCount; i++ {
		res[i] = hexSymbols[rand.IntN(maxRand)]
	}

	return string(res)
}
