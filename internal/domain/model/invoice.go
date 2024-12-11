package model

import "time"

type InvoiceStatus string

const (
	StatusPending    InvoiceStatus = "pending"
	StatusProcessing InvoiceStatus = "processing"
	StatusPaid       InvoiceStatus = "paid"
	StatusError      InvoiceStatus = "error"
)

type Invoice struct {
	ID           uint          // 請求書ID
	Organization *Organization // 請求元企業
	Client       *Client       // 請求先取引先
	IssueDate    time.Time     // 発行日
	Amount       float64       // 支払金額
	Fee          float64       // 手数料
	FeeRate      float64       // 手数料率
	Tax          float64       // 消費税
	TaxRate      float64       // 消費税率
	TotalAmount  float64       // 請求金額
	DueDate      time.Time     // 支払期日
	Status       InvoiceStatus // ステータス
}
