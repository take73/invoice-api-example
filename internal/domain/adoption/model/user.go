package model

type User struct {
	ID             int    // ユーザーID
	OrganizationID int    // 紐づく企業ID
	Name           string // 氏名
	Email          string // メールアドレス
	Password       string // パスワード（ハッシュ化された値を保持）
}
