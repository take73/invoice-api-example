package rdb

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/shopspring/decimal"
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
		input   input
		want    *model.Invoice
		wantErr error
	}{
		{
			name: "1件登録",
			input: input{
				invoice: &model.Invoice{
					Organization: &model.Organization{
						ID:   1,
						Name: "株式会社サンプル",
					},
					Client: &model.Client{
						ID:   1,
						Name: "取引先A",
					},
					IssueDate:   time.Date(2018, 04, 15, 0, 0, 0, 0, time.Local),
					Amount:      decimal.NewFromInt(10000),
					Fee:         decimal.NewFromInt(400),
					FeeRate:     0.04,
					Tax:         decimal.NewFromInt(40),
					TaxRate:     0.1,
					TotalAmount: decimal.NewFromInt(10440),
					DueDate:     time.Date(2018, 04, 30, 0, 0, 0, 0, time.Local),
					Status:      model.StatusPending,
				},
			},
			want: &model.Invoice{
				ID: 1,
				Organization: &model.Organization{
					ID:   1,
					Name: "株式会社サンプル",
				},
				Client: &model.Client{
					ID:   1,
					Name: "取引先A",
				},
				IssueDate:   time.Date(2018, 04, 15, 0, 0, 0, 0, time.Local),
				Amount:      decimal.NewFromInt(10000),
				Fee:         decimal.NewFromInt(400),
				FeeRate:     0.04,
				Tax:         decimal.NewFromInt(40),
				TaxRate:     0.1,
				TotalAmount: decimal.NewFromInt(10440),
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

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("api got != want (-got +want)\n%s", diff)
				return
			}
		})
	}
}

func Test_InvoiceRepository_FindByDueDateRange(t *testing.T) {
	db := testutils.SetupTestDB(testutils.GetFuncName())
	db.Logger = db.Logger.LogMode(logger.Info)
	// テストデータの挿入
	testutils.ExecSQLFile(db, "testdata/test_invoice_repository_find_by_due_date_range.sql")

	type input struct {
		startDate time.Time
		endDate   time.Time
	}

	tests := []struct {
		name    string
		input   input
		want    []*model.Invoice
		wantErr error
	}{
		{
			name: "2件取得",
			input: input{
				startDate: time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC),
				endDate:   time.Date(2024, 2, 2, 23, 59, 59, 0, time.UTC),
			},
			want: []*model.Invoice{
				{
					ID:           2,
					Organization: &model.Organization{ID: 1, Name: "株式会社サンプル"},
					Client:       &model.Client{ID: 2, Name: "取引先B"},
					IssueDate:    time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
					Amount:       decimal.NewFromInt(20000),
					Fee:          decimal.NewFromInt(800),
					FeeRate:      0.04,
					Tax:          decimal.NewFromInt(80),
					TaxRate:      0.1,
					TotalAmount:  decimal.NewFromInt(20880),
					DueDate:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
					Status:       model.StatusProcessing,
				},
				{
					ID:           3,
					Organization: &model.Organization{ID: 2, Name: "有限会社テスト"},
					Client:       &model.Client{ID: 3, Name: "取引先C"},
					IssueDate:    time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
					Amount:       decimal.NewFromInt(30000),
					Fee:          decimal.NewFromInt(1200),
					FeeRate:      0.04,
					Tax:          decimal.NewFromInt(120),
					TaxRate:      0.1,
					TotalAmount:  decimal.NewFromInt(31320),
					DueDate:      time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
					Status:       model.StatusPaid,
				},
				{
					ID:           4,
					Organization: &model.Organization{ID: 2, Name: "有限会社テスト"},
					Client:       &model.Client{ID: 1, Name: "取引先A"},
					IssueDate:    time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
					Amount:       decimal.NewFromInt(40000),
					Fee:          decimal.NewFromInt(1600),
					FeeRate:      0.04,
					Tax:          decimal.NewFromInt(160),
					TaxRate:      0.1,
					TotalAmount:  decimal.NewFromInt(41760),
					DueDate:      time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
					Status:       model.StatusError,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewInvoiceRepository(db)
			got, err := repo.FindByDueDateRange(tt.input.startDate, tt.input.endDate)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("error is nil, want %v", tt.wantErr)
				}
				return
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("api got != want (-got +want)\n%s", diff)
				return
			}
		})
	}
}
