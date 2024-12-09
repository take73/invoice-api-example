package http

import (
	"github.com/labstack/echo/v4"
	"github.com/take73/invoice-api-example/internal/application"
)

func RegisterRoutes(e *echo.Echo, invoiceService *application.InvoiceUsecase) {
	handler := NewInvoiceHandler(invoiceService)

	// ルート設定
	e.POST("/invoice", handler.CreateInvoice)
}
