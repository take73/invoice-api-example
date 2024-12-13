package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/take73/invoice-api-example/internal/infrastructure/http/middleware/testutil"
)

func Test_Middleware_EnsureValidTokenWithScopes(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name               string
		scope              string // middlewareに設定するスコープ
		token              func() (string, error)
		expectedStatusCode int
		expectedMessage    string
	}{
		{
			name:               "read:invoiceをもっているトークン",
			scope:              "read:invoice",
			token:              func() (string, error) { return testutil.NewClientAToken() },
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "Access granted",
		},
		{
			name:               "write:invoiceをもっているトークン",
			scope:              "read:invoice",
			token:              func() (string, error) { return testutil.NewClientBToken() },
			expectedStatusCode: http.StatusForbidden,
			expectedMessage:    "Insufficient scope: required [read:invoice]",
		},
		{
			name:               "read:invoice,write:invoiceどちらもないトークン",
			scope:              "read:invoice",
			token:              func() (string, error) { return testutil.NewClientCToken() },
			expectedStatusCode: http.StatusForbidden,
			expectedMessage:    "Insufficient scope: required [read:invoice]",
		},
		{
			name:               "tokenなし",
			scope:              "read:invoice",
			token:              func() (string, error) { return "", nil },
			expectedStatusCode: http.StatusUnauthorized,
			expectedMessage:    "Missing Authorization header",
		},
		{
			name:               "擬似的になトークン",
			scope:              "read:invoice",
			token:              func() (string, error) { return testutil.GenerateMockJWT() },
			expectedStatusCode: http.StatusUnauthorized,
			expectedMessage:    "Failed to validate token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Middleware to test
			e.Use(EnsureValidTokenWithScopes(tt.scope))

			// Mock handler
			e.GET("/secure-data", func(c echo.Context) error {
				return c.JSON(http.StatusOK, map[string]string{
					"message": "Access granted",
				})
			})

			// Create a new HTTP request with the Authorization header
			req := httptest.NewRequest(http.MethodGet, "/secure-data", nil)

			token, err := tt.token()
			assert.NoError(t, err)

			if token != "" {
				req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
			}

			rec := httptest.NewRecorder()
			e.NewContext(req, rec)

			// Serve the request
			e.ServeHTTP(rec, req)

			// Assert response status code and message
			assert.Equal(t, tt.expectedStatusCode, rec.Code)
			if tt.expectedMessage != "" {
				var response map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedMessage, response["message"])
			}
		})
	}
}
