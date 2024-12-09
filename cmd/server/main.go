package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/take73/invoice-api-example/internal/application"
	"github.com/take73/invoice-api-example/internal/infrastructure/http"
	"github.com/take73/invoice-api-example/internal/infrastructure/rdb"
)

func main() {
	// 依存関係の組み立て
	invoiceRepo := rdb.NewInvoiceRepository()
	invoiceService := application.NewInvoiceService(invoiceRepo)

	// Echoルータの初期化
	e := echo.New()
	http.RegisterRoutes(e, invoiceService)

	// サーバーの起動
	log.Println("Starting server on :8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
