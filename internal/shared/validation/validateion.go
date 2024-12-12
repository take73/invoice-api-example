package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/take73/invoice-api-example/internal/shared/types"
)

// CustomValidator Validatorの構造体
type CustomValidator struct {
	Validator *validator.Validate
}

// NewCustomValidator コンストラクタ
func NewCustomValidator() *CustomValidator {
	v := validator.New()
	// カスタムバリデーションの登録
	v.RegisterValidation("required_custom_date", validateCustomDate)
	return &CustomValidator{Validator: v}
}

// Validate 構造体のフィールドを検証する
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

// validateCustomDate カスタム日付型のゼロ値を検出
func validateCustomDate(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(types.CustomDate)
	if !ok {
		return false
	}
	return !date.Time.IsZero() // IsZeroだったら失敗
}

// ValidRate 割合の妥当性（0.0から1.0の範囲内かどうか）を検証します
// TODO: validator と統合できるかも
func ValidRate(rate float64) bool {
	return rate >= 0.0 && rate <= 1.0
}
