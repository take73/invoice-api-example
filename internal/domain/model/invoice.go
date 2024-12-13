package model

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

type InvoiceStatus string

const (
	StatusPending    InvoiceStatus = "pending"
	StatusProcessing InvoiceStatus = "processing"
	StatusPaid       InvoiceStatus = "paid"
	StatusError      InvoiceStatus = "error"
)

type Invoice struct {
	ID           uint            // 請求書ID
	Organization *Organization   // 請求元企業
	Client       *Client         // 請求先取引先
	IssueDate    time.Time       // 発行日
	Amount       decimal.Decimal // 支払金額
	Fee          decimal.Decimal // 手数料
	FeeRate      float64         // 手数料率
	Tax          decimal.Decimal // 消費税
	TaxRate      float64         // 消費税率
	TotalAmount  decimal.Decimal // 請求金額
	DueDate      time.Time       // 支払期日
	Status       InvoiceStatus   // ステータス
}

const defaultFeeRate = 0.04

func NewInvoice(org *Organization, client *Client, amount int64, issueDate, dueDate time.Time) (*Invoice, error) {
	feeRate := defaultFeeRate
	if feeRateStr := os.Getenv("FEE_RATE"); feeRateStr != "" {
		if parsedFeeRate, err := strconv.ParseFloat(feeRateStr, 64); err == nil {
			feeRate = parsedFeeRate
		} else {
			return nil, fmt.Errorf("invalid FEE_RATE value: %v", err)
		}
	}

	return &Invoice{
		Organization: org,
		Client:       client,
		Amount:       decimal.NewFromInt(amount),
		FeeRate:      feeRate,
		IssueDate:    issueDate,
		DueDate:      dueDate,
		Status:       StatusPending,
	}, nil
}

// Calculate 手数料、消費税、請求金額を計算してセットする
func (i *Invoice) Calculate(taxRate float64) {
	// 支払金額 (Amount) を Decimal に変換
	amount := i.Amount
	feeRate := decimal.NewFromFloat(i.FeeRate)
	taxRateDecimal := decimal.NewFromFloat(taxRate)

	// 手数料を計算: Fee = Amount * FeeRate
	fee := amount.Mul(feeRate)
	i.Fee = fee

	// 消費税を計算: Tax = Fee * TaxRate
	tax := fee.Mul(taxRateDecimal)
	i.Tax = tax

	// 請求金額を計算: TotalAmount = Amount + Fee + Tax
	totalAmount := amount.Add(fee).Add(tax)
	i.TotalAmount = totalAmount

	// 消費税率をセット
	i.TaxRate = taxRate
}

// truncateDecimalToInt 小数点以下を切り捨てて int で返す
func truncateDecimalToInt(d decimal.Decimal) int64 {
	// 小数点以下を切り捨てる
	truncated := d.Truncate(0)
	return int64(truncated.IntPart())
}

// AmountAsInt 小数点以下を切り捨てて int で返す
func (i *Invoice) AmountAsInt() int64 {
	return truncateDecimalToInt(i.Amount)
}

// TotalAmountAsInt 小数点以下を切り捨てて int で返す
func (i *Invoice) TotalAmountAsInt() int64 {
	return truncateDecimalToInt(i.TotalAmount)
}

// TaxAsInt 小数点以下を切り捨てて int で返す
func (i *Invoice) TaxAsInt() int64 {
	return truncateDecimalToInt(i.Tax)
}

// FeeAsInt 小数点以下を切り捨てて int で返す
func (i *Invoice) FeeAsInt() int64 {
	return truncateDecimalToInt(i.Fee)
}
