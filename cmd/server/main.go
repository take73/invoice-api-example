package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/take73/invoice-api-example/internal/application"
	myHttp "github.com/take73/invoice-api-example/internal/infrastructure/http"
	"github.com/take73/invoice-api-example/internal/infrastructure/rdb"
	"github.com/take73/invoice-api-example/internal/shared/validation"
	"gorm.io/gorm/logger"
)

func main() {
	db, err := rdb.NewDB()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
		return
	}
	db.Logger = db.Logger.LogMode(logger.Info)

	// 依存関係を設定、肥大化したらwire等のDIの仕組みを導入する
	invoiceRepo := rdb.NewInvoiceRepository(db)
	clientRepo := rdb.NewClientRepository(db)
	organizationRepo := rdb.NewOrganizationRepository(db)
	taxRateRepo := rdb.NewTaxRateRepository(db)
	invoiceUsecase := application.NewInvoiceUsecase(invoiceRepo, clientRepo, organizationRepo, taxRateRepo)

	e := echo.New()
	e.Validator = validation.NewCustomValidator()
	myHttp.RegisterRoutes(e, invoiceUsecase)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Start server
	go func() {
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second) // 20秒待つ
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
