package rdb

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/take73/invoice-api-example/internal/infrastructure/rdb/testutils"
	"gorm.io/gorm/logger"
)

func Test_TaxRateRepository_GetRateByDate(t *testing.T) {
	db := testutils.SetupTestDB(testutils.GetFuncName())
	db.Logger = db.Logger.LogMode(logger.Info)

	type input struct {
		date time.Time
	}

	tests := []struct {
		name    string
		before  func()
		input   input
		want    float64
		wantErr error
	}{
		{
			name: "1件取得",
			input: input{
				date: time.Date(2021, 4, 1, 0, 0, 0, 0, time.Local),
			},
			want: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewTaxRateRepository(db)
			got, err := repo.GetRateByDate(tt.input.date)

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
