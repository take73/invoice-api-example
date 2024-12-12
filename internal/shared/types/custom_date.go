package types

import "time"

type CustomDate struct {
	time.Time
}

const dateFormat = "2006-01-02"

// / UnmarshalJSON JSONフィールドから日付をデコード
func (d *CustomDate) UnmarshalJSON(b []byte) error {
	str := string(b)
	// JSON の場合、クオートで囲まれているので削除
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}
	return d.unmarshalCommon(str)
}

// UnmarshalParam クエリパラメータから日付をデコード
func (d *CustomDate) UnmarshalParam(param string) error {
	return d.unmarshalCommon(param)
}

// unmarshalCommon 実際のパース処理を共通化
func (d *CustomDate) unmarshalCommon(dateStr string) error {
	if dateStr == "" {
		return nil // 空文字の場合はスキップ
	}
	parsedTime, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		return err
	}
	d.Time = parsedTime
	return nil
}

// MarshalJSON 日付をJSON形式でエンコード
func (d CustomDate) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil // ゼロ値の場合は null を返す
	}
	return []byte(`"` + d.Time.Format(dateFormat) + `"`), nil
}
