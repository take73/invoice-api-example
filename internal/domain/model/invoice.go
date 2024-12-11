package model

import (
	"fmt"
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
