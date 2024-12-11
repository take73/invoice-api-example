package testutils

import (
	"math/big"
)

// CompareBigRat 比較用カスタム関数
func CompareBigRat(x, y *big.Rat) bool {
	if x == nil || y == nil {
		return x == y
	}
	return x.Cmp(y) == 0
}
