package http

import (
	"github.com/labstack/echo/v4"
	"github.com/take73/invoice-api-example/internal/application"
	"github.com/take73/invoice-api-example/internal/infrastructure/http/middleware"
)

func RegisterRoutes(e *echo.Echo, invoiceUsecase application.InvoiceUsecase) {
	handler := NewInvoiceHandler(invoiceUsecase)

	// ルート設定
	e.POST("/invoice", handler.CreateInvoice, middleware.AuthWithScopes("write:invoice"))
	e.GET("/invoice", handler.ListInvoice, middleware.AuthWithScopes("read:invoice"))
}
