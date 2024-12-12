package model_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/shopspring/decimal"
	"github.com/take73/invoice-api-example/internal/domain/model"
)
func Test_Invoice_Calculate(t *testing.T) {
	tests := []struct {
		name    string
		amount  int64
		feeRate float64
		taxRate float64
		want    model.Invoice
	}{
		{
			name:    "正常ケース: 手数料4%, 消費税10%",
			amount:  10000,
			feeRate: 0.04,
			taxRate: 0.1,
			want: model.Invoice{
				Amount:      decimal.NewFromInt(10000),
				FeeRate:     0.04,
				Fee:         decimal.NewFromInt(400),
				Tax:         decimal.NewFromInt(40),
				TaxRate:     0.1,
				TotalAmount: decimal.NewFromInt(10440),
			},
		},
		{
			name:    "非常に大きな金額: 手数料4%, 消費税10%",
			amount:  1_000_000_000_000, // 1兆円
			feeRate: 0.04,
			taxRate: 0.1,
			want: model.Invoice{
				Amount:      decimal.NewFromInt(1_000_000_000_000),
				FeeRate:     0.04,
				Fee:         decimal.NewFromInt(40_000_000_000),
				Tax:         decimal.NewFromInt(4_000_000_000),
				TaxRate:     0.1,
				TotalAmount: decimal.NewFromInt(1_044_000_000_000),
			},
		},
		{
			name:    "手数料率が0: 消費税10%",
			amount:  5000,
			feeRate: 0.0, // 手数料率0
			taxRate: 0.1,
			want: model.Invoice{
				Amount:      decimal.NewFromInt(5000),
				FeeRate:     0.0,
				Fee:         decimal.NewFromInt(0),
				Tax:         decimal.NewFromInt(0),
				TaxRate:     0.1,
				TotalAmount: decimal.NewFromInt(5000),
			},
		},
		{
			name:    "消費税率が0: 手数料4%",
			amount:  5000,
			feeRate: 0.04,
			taxRate: 0.0, // 消費税率0
			want: model.Invoice{
				Amount:      decimal.NewFromInt(5000),
				FeeRate:     0.04,
				Fee:         decimal.NewFromInt(200),
				Tax:         decimal.NewFromInt(0),
				TaxRate:     0.0,
				TotalAmount: decimal.NewFromInt(5200),
			},
		},
		{
			name:    "手数料率と消費税率がともに0",
			amount:  10000,
			feeRate: 0.0, // 手数料率0
			taxRate: 0.0, // 消費税率0
			want: model.Invoice{
				Amount:      decimal.NewFromInt(10000),
				FeeRate:     0.0,
				Fee:         decimal.NewFromInt(0),
				Tax:         decimal.NewFromInt(0),
				TaxRate:     0.0,
				TotalAmount: decimal.NewFromInt(10000),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invoice := model.Invoice{
				Amount:  decimal.NewFromInt(tt.amount),
				FeeRate: tt.feeRate,
			}

			invoice.Calculate(tt.taxRate)

			// 結果を比較
			if diff := cmp.Diff(tt.want.Fee, invoice.Fee); diff != "" {
				t.Errorf("Fee mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.want.Tax, invoice.Tax); diff != "" {
				t.Errorf("Tax mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.want.TotalAmount, invoice.TotalAmount); diff != "" {
				t.Errorf("TotalAmount mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_Invoice_TotalAmountAsInt(t *testing.T) {
	tests := []struct {
		name        string
		totalAmount decimal.Decimal
		expectedInt int64
		expectError bool
	}{
		{
			name:        "正常値 (整数値)",
			totalAmount: decimal.NewFromInt(10440),
			expectedInt: 10440,
			expectError: false,
		},
		{
			name:        "小数点以下切り捨て (10.99 -> 10)",
			totalAmount: decimal.NewFromFloat(10.99),
			expectedInt: 10,
			expectError: false,
		},
		{
			name:        "小数点以下切り捨て (10.00000001 -> 10)",
			totalAmount: decimal.NewFromFloat(10.00000001),
			expectedInt: 10,
			expectError: false,
		},
		{
			name:        "ゼロ (0.999 -> 0)",
			totalAmount: decimal.NewFromFloat(0.999),
			expectedInt: 0,
			expectError: false,
		},
		{
			name:        "負の値 (-10.99 -> -10)",
			totalAmount: decimal.NewFromFloat(-10.99),
			expectedInt: -10,
			expectError: false,
		},
		{
			name:        "非常に小さい値",
			totalAmount: decimal.NewFromFloat(0.00000001),
			expectedInt: 0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Invoice オブジェクトを作成
			invoice := &model.Invoice{
				TotalAmount: tt.totalAmount,
			}

			// TotalAmountAsInt をテスト
			got := invoice.TotalAmountAsInt()

			// 結果の検証
			if got != tt.expectedInt {
				t.Errorf("unexpected result: got %d, want %d", got, tt.expectedInt)
			}
		})
	}
}
