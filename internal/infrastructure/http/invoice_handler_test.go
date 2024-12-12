package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/take73/invoice-api-example/internal/application"
	"github.com/take73/invoice-api-example/internal/infrastructure/http/testutils"
	commonErrors "github.com/take73/invoice-api-example/internal/shared/errors"
	"github.com/take73/invoice-api-example/internal/shared/validation"
)

func Test_InvoiceHandler_CreateInvoice(t *testing.T) {
	e := echo.New()
	e.Validator = validation.NewCustomValidator()

	tests := []struct {
		name           string
		setupMock      func(*testutils.MockInvoiceUsecase)
		payload        map[string]interface{}
		expectedStatus int
		expectedBody   func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			setupMock: func(mockUsecase *testutils.MockInvoiceUsecase) {
				mockUsecase.On("CreateInvoice", application.CreateInvoiceDto{
					UserID:    1,
					ClientID:  1,
					IssueDate: time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC),
					Amount:    10000,
					DueDate:   time.Date(2023, 12, 15, 0, 0, 0, 0, time.UTC),
				}).Return(&application.CreatedInvoiceDto{
					ID:               1,
					OrganizationID:   1,
					OrganizationName: "Test Organization",
					ClientID:         1,
					ClientName:       "Test Client",
					IssueDate:        time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC),
					Amount:           10000,
					Fee:              400,
					FeeRate:          0.04,
					Tax:              40,
					TaxRate:          0.1,
					TotalAmount:      10440,
					DueDate:          time.Date(2023, 12, 15, 0, 0, 0, 0, time.UTC),
					Status:           "pending",
				}, nil)
			},
			payload: map[string]interface{}{
				"userId":    1,
				"clientId":  1,
				"issueDate": "2023-12-01",
				"amount":    10000,
				"dueDate":   "2023-12-15",
			},
			expectedStatus: http.StatusOK,
			expectedBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response CreateInvoiceResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, uint(1), response.ID)
				assert.Equal(t, "Test Organization", response.OrganizationName)
				assert.Equal(t, "Test Client", response.ClientName)
			},
		},
		{
			name:      "userIdが-1の場合, invalid request",
			setupMock: func(mockUsecase *testutils.MockInvoiceUsecase) {}, // Mock is not called in this case
			payload: map[string]interface{}{
				"userId":    -1,
				"clientId":  1,
				"issueDate": "2023-12-01",
				"amount":    10000,
				"dueDate":   "2023-12-15",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "invalid request", response["error"])
			},
		},
		{
			name:      "userIdが0の場合, validation failed",
			setupMock: func(mockUsecase *testutils.MockInvoiceUsecase) {}, // Mock is not called in this case
			payload: map[string]interface{}{
				"userId":    0,
				"clientId":  1,
				"issueDate": "2023-12-01",
				"amount":    10000,
				"dueDate":   "2023-12-15",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "validation failed", response["error"])
			},
		},
		{
			name:      "clientIdが-1の場合, invalid request",
			setupMock: func(mockUsecase *testutils.MockInvoiceUsecase) {}, // Mock is not called in this case
			payload: map[string]interface{}{
				"userId":    1,
				"clientId":  -1,
				"issueDate": "2023-12-01",
				"amount":    10000,
				"dueDate":   "2023-12-15",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "invalid request", response["error"])
			},
		},
		{
			name:      "clientIdが0の場合, validation failed",
			setupMock: func(mockUsecase *testutils.MockInvoiceUsecase) {}, // Mock is not called in this case
			payload: map[string]interface{}{
				"userId":    1,
				"clientId":  0,
				"issueDate": "2023-12-01",
				"amount":    10000,
				"dueDate":   "2023-12-15",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "validation failed", response["error"])
			},
		},
		{
			name:      "issueDateのformatが不正の場合, invalid request",
			setupMock: func(mockUsecase *testutils.MockInvoiceUsecase) {}, // Mock is not called in this case
			payload: map[string]interface{}{
				"userId":    1,
				"clientId":  1,
				"issueDate": "2023/12/01",
				"amount":    10000,
				"dueDate":   "2023-12-15",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "invalid request", response["error"])
			},
		},
		{
			name:      "issueDateがない場合, validation failed",
			setupMock: func(mockUsecase *testutils.MockInvoiceUsecase) {}, // Mock is not called in this case
			payload: map[string]interface{}{
				"userId":   1,
				"clientId": 1,
				"amount":   10000,
				"dueDate":  "2023-12-15",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "validation failed", response["error"])
			},
		},
		{
			name: "amountが0の場合, success",
			setupMock: func(mockUsecase *testutils.MockInvoiceUsecase) {
				mockUsecase.On("CreateInvoice", application.CreateInvoiceDto{
					UserID:    1,
					ClientID:  1,
					IssueDate: time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC),
					Amount:    0,
					DueDate:   time.Date(2023, 12, 15, 0, 0, 0, 0, time.UTC),
				}).Return(&application.CreatedInvoiceDto{
					ID:               1,
					OrganizationID:   1,
					OrganizationName: "Test Organization",
					ClientID:         1,
					ClientName:       "Test Client",
					IssueDate:        time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC),
					Amount:           0,
					Fee:              0,
					FeeRate:          0.04,
					Tax:              0,
					TaxRate:          0.1,
					TotalAmount:      10440,
					DueDate:          time.Date(2023, 12, 15, 0, 0, 0, 0, time.UTC),
					Status:           "pending",
				}, nil)
			},
			payload: map[string]interface{}{
				"userId":    1,
				"clientId":  1,
				"issueDate": "2023-12-01",
				"amount":    0,
				"dueDate":   "2023-12-15",
			},
			expectedStatus: http.StatusOK,
			expectedBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response CreateInvoiceResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, uint(1), response.ID)
				assert.Equal(t, "Test Organization", response.OrganizationName)
				assert.Equal(t, "Test Client", response.ClientName)
			},
		},
		{
			name:      "dueDateのformatが不正の場合, invalid request",
			setupMock: func(mockUsecase *testutils.MockInvoiceUsecase) {}, // Mock is not called in this case
			payload: map[string]interface{}{
				"userId":    1,
				"clientId":  1,
				"issueDate": "2023-12-01",
				"amount":    10000,
				"dueDate":   "2023/12/15",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "invalid request", response["error"])
			},
		},
		{
			name:      "dueDateがない場合, validation failed",
			setupMock: func(mockUsecase *testutils.MockInvoiceUsecase) {}, // Mock is not called in this case
			payload: map[string]interface{}{
				"userId":    1,
				"clientId":  1,
				"issueDate": "2023-12-01",
				"amount":    10000,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "validation failed", response["error"])
			},
		},
		{
			name: "usecaseでNoFoundエラーが発生した場合, related company or client not found",
			setupMock: func(mockUsecase *testutils.MockInvoiceUsecase) {
				mockUsecase.On("CreateInvoice", application.CreateInvoiceDto{
					UserID:    1,
					ClientID:  1,
					IssueDate: time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC),
					Amount:    10000,
					DueDate:   time.Date(2023, 12, 15, 0, 0, 0, 0, time.UTC),
				}).Return(nil, commonErrors.ErrNotFound)
			},
			payload: map[string]interface{}{
				"userId":    1,
				"clientId":  1,
				"issueDate": "2023-12-01",
				"amount":    10000,
				"dueDate":   "2023-12-15",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "related company or client not found", response["error"])
			},
		},
		{
			name: "usecaseでエラーが発生した場合, could not create invoice",
			setupMock: func(mockUsecase *testutils.MockInvoiceUsecase) {
				mockUsecase.On("CreateInvoice", application.CreateInvoiceDto{
					UserID:    1,
					ClientID:  1,
					IssueDate: time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC),
					Amount:    10000,
					DueDate:   time.Date(2023, 12, 15, 0, 0, 0, 0, time.UTC),
				}).Return(nil, errors.New("unexpected error"))
			},
			payload: map[string]interface{}{
				"userId":    1,
				"clientId":  1,
				"issueDate": "2023-12-01",
				"amount":    10000,
				"dueDate":   "2023-12-15",
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "could not create invoice", response["error"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 新しいモックインスタンスを作成
			mockUsecase := &testutils.MockInvoiceUsecase{}
			tt.setupMock(mockUsecase)

			// ハンドラを新規作成
			handler := NewInvoiceHandler(mockUsecase)

			// リクエストのセットアップ
			reqBody, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/invoice", bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// ハンドラの実行
			err := handler.CreateInvoice(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			// レスポンスの検証
			if tt.expectedBody != nil {
				tt.expectedBody(t, rec)
			}

			// モックのアサーション
			mockUsecase.AssertExpectations(t)
		})
	}
}
