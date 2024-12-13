package rdb

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/take73/invoice-api-example/internal/domain/model"
	"github.com/take73/invoice-api-example/internal/infrastructure/rdb/testutils"
	"gorm.io/gorm/logger"
)

func Test_OrganizationRepository_GetByID(t *testing.T) {
	db, cleanup := testutils.SetupTestDB(testutils.GetFuncName())
	defer cleanup()

	db.Logger = db.Logger.LogMode(logger.Info)
	// testutils.ExecSQLFile(db, "testdata/test_organization_reposiroty_get_by_id.sql")

	type input struct {
		id uint
	}

	tests := []struct {
		name    string
		before  func()
		input   input
		want    *model.Organization
		wantErr error
	}{
		{
			name: "1件取得",
			input: input{
				id: 1,
			},
			want: &model.Organization{
				ID:             1,
				Name:           "株式会社サンプル",
				Representative: "山田 太郎",
				PhoneNumber:    "03-1234-5678",
				PostalCode:     "100-0001",
				Address:        "東京都千代田区丸の内1-1-1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewOrganizationRepository(db)
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

func Test_OrganizationRepository_GetByUserID(t *testing.T) {
	db, cleanup := testutils.SetupTestDB(testutils.GetFuncName())
	defer cleanup()

	db.Logger = db.Logger.LogMode(logger.Info)
	// testutils.ExecSQLFile(db, "testdata/test_organization_reposiroty_get_by_user_id.sql")

	type input struct {
		userID uint
	}

	tests := []struct {
		name    string
		before  func()
		input   input
		want    *model.Organization
		wantErr error
	}{
		{
			name: "1件取得",
			input: input{
				userID: 1,
			},
			want: &model.Organization{
				ID:             1,
				Name:           "株式会社サンプル",
				Representative: "山田 太郎",
				PhoneNumber:    "03-1234-5678",
				PostalCode:     "100-0001",
				Address:        "東京都千代田区丸の内1-1-1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewOrganizationRepository(db)
			got, err := repo.GetByUserID(tt.input.userID)

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
