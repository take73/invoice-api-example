package entity

import (
	"database/sql/driver"
	"fmt"
	"math/big"
)

// BigRat is a wrapper for big.Rat to implement GORM interfaces
type BigRat struct {
	*big.Rat
}

// Scan データベースから読み込む際の処理
func (b *BigRat) Scan(value interface{}) error {
	if str, ok := value.(string); ok {
		b.Rat = new(big.Rat)
		if _, ok := b.Rat.SetString(str); !ok {
			return fmt.Errorf("failed to parse %s as big.Rat", str)
		}
		return nil
	}
	return fmt.Errorf("failed to scan type %T into BigRat", value)
}

// Value データベースに保存する際の処理
func (b BigRat) Value() (driver.Value, error) {
	if b.Rat == nil {
		return nil, nil
	}
	return b.Rat.FloatString(2), nil // 小数点以下2桁までで保存
}
