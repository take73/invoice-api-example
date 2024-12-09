package http

import (
	"net/http"

	"github.com/take73/invoice-api-example/internal/application"

	"github.com/labstack/echo/v4"
)

type InvoiceHandler struct {
	service *application.InvoiceUsecase
}

func NewInvoiceHandler(service *application.InvoiceUsecase) *InvoiceHandler {
	return &InvoiceHandler{service: service}
}

func (h *InvoiceHandler) CreateInvoice(c echo.Context) error {
	// リクエストボディのパース（ここでは単純な例として空データを渡します）
	var data interface{}
	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	// サービスロジック呼び出し
	if err := h.service.CreateInvoice(data); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not create invoice"})
	}

	// 成功時のレスポンス
	return c.JSON(http.StatusOK, map[string]string{"message": "invoice created"})
}
