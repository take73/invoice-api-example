package http

import (
	"github.com/labstack/echo/v4"
	"github.com/take73/invoice-api-example/internal/application"
)

func RegisterRoutes(e *echo.Echo, invoiceUsecase *application.InvoiceUsecase) {
	handler := NewInvoiceHandler(invoiceUsecase)

	// ルート設定
	e.POST("/invoice", handler.CreateInvoice)
	e.GET("/invoice", handler.ListInvoice)
}
