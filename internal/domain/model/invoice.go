package model

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"os"
	"strconv"
	"time"
)

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
	Amount       *big.Rat      // 支払金額
	Fee          *big.Rat      // 手数料
	FeeRate      float64       // 手数料率
	Tax          *big.Rat      // 消費税
	TaxRate      float64       // 消費税率
	TotalAmount  *big.Rat      // 請求金額
	DueDate      time.Time     // 支払期日
	Status       InvoiceStatus // ステータス
}

const defaultFeeRate = 0.04

func NewInvoice(org *Organization, client *Client, amount float64, issueDate, dueDate time.Time) (*Invoice, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}

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
		Amount:       big.NewRat(0, 1).SetFloat64(amount),
		FeeRate:      feeRate,
		IssueDate:    issueDate,
		DueDate:      dueDate,
		Status:       StatusPending,
	}, nil
}

// Calculate 手数料、消費税、請求金額を計算してセットする
func (i *Invoice) Calculate(taxRate float64) {
	// 手数料を計算: Fee = Amount * FeeRate
	fee := big.NewRat(0, 1).Mul(i.Amount, big.NewRat(0, 1).SetFloat64(i.FeeRate))
	i.Fee = fee

	// 消費税を計算: Tax = Fee * TaxRate
	tax := big.NewRat(0, 1).Mul(fee, big.NewRat(0, 1).SetFloat64(taxRate))
	i.Tax = tax

	// 請求金額を計算: TotalAmount = Amount + Fee + Tax
	totalAmount := big.NewRat(0, 1).Add(i.Amount, fee)
	totalAmount.Add(totalAmount, tax)
	i.TotalAmount = totalAmount

	// 消費税率をセット
	i.TaxRate = taxRate
}

// truncateRatToInt 小数点以下を切り捨てて int で返す
func truncateRatToInt(r *big.Rat) (int, error) {
	value, ok := r.Float64()
	if !ok {
		return 0, errors.New("failed to convert Rat to float64")
	}
	truncated := math.Floor(value)
	return int(truncated), nil
}

// AmountAsInt 小数点以下を切り捨てて intで返す
func (i *Invoice) AmountAsInt() (int, error) {
	return truncateRatToInt(i.Amount)
}

// TotalAmountAsInt 小数点以下を切り捨てて intで返す
func (i *Invoice) TotalAmountAsInt() (int, error) {
	return truncateRatToInt(i.TotalAmount)
}

// TaxAsInt 小数点以下を切り捨てて intで返す
func (i *Invoice) TaxAsInt() (int, error) {
	return truncateRatToInt(i.Tax)
}

// FeeAsInt 小数点以下を切り捨てて intで返す
func (i *Invoice) FeeAsInt() (int, error) {
	return truncateRatToInt(i.Fee)
}
