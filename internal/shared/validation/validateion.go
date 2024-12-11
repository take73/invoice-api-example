package validation

// ValidRate 割合の妥当性（0.0から1.0の範囲内かどうか）を検証します
func ValidRate(rate float64) bool {
	return rate < 0.0 || rate > 1.0
}
