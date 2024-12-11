package model_test

import (
	"math/big"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/take73/invoice-api-example/internal/domain/model"
)

func Test_Invoice_Calculate(t *testing.T) {
	tests := []struct {
		name    string
		amount  float64
		feeRate float64
		taxRate float64
		want    model.Invoice
	}{
		{
			name:    "正常ケース: 手数料4%, 消費税10%",
			amount:  10000.0,
			feeRate: 0.04,
			taxRate: 0.1,
			want: model.Invoice{
				Fee:         big.NewRat(40, 1),    // 10000 * 0.04 = 400
				Tax:         big.NewRat(4, 1),     // 400 * 0.1 = 40
				TotalAmount: big.NewRat(10440, 1), // 10000 + 400 + 40 = 10440
				TaxRate:     0.1,
			},
		},
		{
			name:    "正常ケース: 手数料3%, 消費税8%",
			amount:  5000.0,
			feeRate: 0.03,
			taxRate: 0.08,
			want: model.Invoice{
				Fee:         big.NewRat(150, 1),  // 5000 * 0.03 = 150
				Tax:         big.NewRat(12, 1),   // 150 * 0.08 = 12
				TotalAmount: big.NewRat(5150, 1), // 5000 + 150 + 12 = 5150
				TaxRate:     0.8,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invoice := model.Invoice{
				Amount:  big.NewRat(0, 1).SetFloat64(tt.amount),
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
		totalAmount *big.Rat
		expectedInt int
		expectError bool
	}{
		{
			name:        "正常値 (整数値)",
			totalAmount: big.NewRat(10440, 1), // 10440.0
			expectedInt: 10440,
			expectError: false,
		},
		{
			name:        "小数点以下切り捨て (10.99 -> 10)",
			totalAmount: big.NewRat(1099, 100), // 10.99
			expectedInt: 10,
			expectError: false,
		},
		{
			name:        "小数点以下切り捨て (10.00000001 -> 10)",
			totalAmount: big.NewRat(1000000001, 100000000), // 10.00000001
			expectedInt: 10,
			expectError: false,
		},
		{
			name:        "ゼロ (0.999 -> 0)",
			totalAmount: big.NewRat(999, 1000), // 0.999
			expectedInt: 0,
			expectError: false,
		},
		{
			name:        "負の値 (-10.99 -> -10)",
			totalAmount: big.NewRat(-1099, 100), // -10.99
			expectedInt: -10,
			expectError: false,
		},
		{
			name:        "非常に小さい値",
			totalAmount: big.NewRat(1, 100000000), // 0.00000001
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
			got, err := invoice.TotalAmountAsInt()

			// エラーの検証
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// 結果の検証
			if got != tt.expectedInt {
				t.Errorf("unexpected result: got %d, want %d", got, tt.expectedInt)
			}
		})
	}
}