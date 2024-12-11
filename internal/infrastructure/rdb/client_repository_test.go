package rdb

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/take73/invoice-api-example/internal/domain/model"
	"github.com/take73/invoice-api-example/internal/infrastructure/rdb/testutils"
	"gorm.io/gorm/logger"
)

func Test_ClientRepository_GetByID(t *testing.T) {
	db := testutils.SetupTestDB(testutils.GetFuncName())
	db.Logger = db.Logger.LogMode(logger.Info)
	testutils.ExecSQLFile(db, "testdata/test_client_reposiroty_get_by_id.sql")

	type input struct {
		id uint
	}

	tests := []struct {
		name    string
		before  func()
		input   input
		want    *model.Client
		wantErr error
	}{
		{
			name: "1件登録",
			input: input{
				id: 1,
			},
			want: &model.Client{
				ID:             1,
				OrganizationID: 1,
				Name:           "取引先A",
				Representative: "取引先担当者A",
				PhoneNumber:    "03-1234-0001",
				PostalCode:     "100-0010",
				Address:        "東京都港区芝公園1-1-1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewClientRepository(db)
			got, err := repo.GetByID(tt.input.id)

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
