package model

type Client struct {
	ID             uint   // クライアントID
	OrganizationID uint   // 紐づく組織ID
	Name           string // 法人名
	Representative string // 代表者名
	PhoneNumber    string // 電話番号
	PostalCode     string // 郵便番号
	Address        string // 住所
}
