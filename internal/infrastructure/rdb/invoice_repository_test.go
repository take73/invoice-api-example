package rdb

import (
	"math/big"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/take73/invoice-api-example/internal/domain/model"
	"github.com/take73/invoice-api-example/internal/infrastructure/rdb/testutils"
	"gorm.io/gorm/logger"
)

func Test_InvoiceRepository_Create(t *testing.T) {
	db := testutils.SetupTestDB(testutils.GetFuncName())
	db.Logger = db.Logger.LogMode(logger.Info)

	type input struct {
		invoice *model.Invoice
	}

	tests := []struct {
		name    string
		before  func()
		input   input
		want    *model.Invoice
		wantErr error
	}{
		{
			name: "1件登録",
			input: input{
				invoice: &model.Invoice{
					Organization: &model.Organization{
						ID: 1,
					},
					Client: &model.Client{
						ID: 1,
					},
					IssueDate:   time.Date(2018, 04, 15, 0, 0, 0, 0, time.Local),
					Amount:      big.NewRat(10000, 1),
					Fee:         big.NewRat(400, 1),
					FeeRate:     0.04,
					Tax:         big.NewRat(1000, 1),
					TaxRate:     0.1,
					TotalAmount: big.NewRat(10440, 1),
					DueDate:     time.Date(2018, 04, 30, 0, 0, 0, 0, time.Local),
					Status:      model.StatusPending,
				},
			},
			want: &model.Invoice{
				ID: 1,
				Organization: &model.Organization{
					ID: 1,
				},
				Client: &model.Client{
					ID: 1,
				},
				IssueDate:   time.Date(2018, 04, 15, 0, 0, 0, 0, time.Local),
				Amount:      big.NewRat(10000, 1),
				Fee:         big.NewRat(400, 1),
				FeeRate:     0.04,
				Tax:         big.NewRat(1000, 1),
				TaxRate:     0.1,
				TotalAmount: big.NewRat(10440, 1),
				DueDate:     time.Date(2018, 04, 30, 0, 0, 0, 0, time.Local),
				Status:      model.StatusPending,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewInvoiceRepository(db)
			got, err := repo.Create(tt.input.invoice)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("error is nil, want %v", tt.wantErr)
				}
				return
			}

			opts := cmp.Options{
				cmp.Comparer(testutils.CompareBigRat),
			}

			if diff := cmp.Diff(got, tt.want, opts); diff != "" {
				t.Errorf("api got != want (-got +want)\n%s", diff)
				return
			}
		})
	}

}
