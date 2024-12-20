package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/take73/invoice-api-example/internal/application"
	"github.com/take73/invoice-api-example/internal/infrastructure/http"
	"github.com/take73/invoice-api-example/internal/infrastructure/rdb"
	"github.com/take73/invoice-api-example/internal/shared/validation"
)

func main() {
	db, err := rdb.NewDB()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
		return
	}

	// 依存関係を設定、肥大化したらwire等のDIの仕組みを導入する
	invoiceRepo := rdb.NewInvoiceRepository(db)
	clientRepo := rdb.NewClientRepository(db)
	organizationRepo := rdb.NewOrganizationRepository(db)
	taxRateRepo := rdb.NewTaxRateRepository(db)
	invoiceUsecase := application.NewInvoiceUsecase(invoiceRepo, clientRepo, organizationRepo, taxRateRepo)

	e := echo.New()
	e.Validator = validation.NewCustomValidator()
	http.RegisterRoutes(e, invoiceUsecase)

	log.Println("Starting server on :1323")
	if err := e.Start(":1323"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
