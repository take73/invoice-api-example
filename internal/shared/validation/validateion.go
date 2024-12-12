package validation

import (
	"github.com/go-playground/validator/v10"
)

// CustomValidator Validatorの構造体
type CustomValidator struct {
	Validator *validator.Validate
}

// NewCustomValidator コンストラクタ
func NewCustomValidator() *CustomValidator {
	v := validator.New()
	// カスタムバリデーションの登録
	// v.RegisterValidation("date", validateDate)
	return &CustomValidator{Validator: v}
}

// Validate 構造体のフィールドを検証する
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

// ValidRate 割合の妥当性（0.0から1.0の範囲内かどうか）を検証します
// TODO: validator と統合できるかも
func ValidRate(rate float64) bool {
	return rate >= 0.0 && rate <= 1.0
}
