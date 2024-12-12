package http

import (
	"log"
	"net/http"
	"time"

	"github.com/take73/invoice-api-example/internal/application"

	"github.com/labstack/echo/v4"
)

type InvoiceHandler struct {
	usecase *application.InvoiceUsecase
}

func NewInvoiceHandler(service *application.InvoiceUsecase) *InvoiceHandler {
	return &InvoiceHandler{usecase: service}
}

type CreateInvoiceRequest struct {
	UserID    uint       `json:"userId"`
	ClientID  uint       `json:"clientId"`
	IssueDate CustomDate `json:"issueDate"`
	Amount    int64      `json:"amount"`
	DueDate   CustomDate `json:"dueDate"`
}

type CreateInvoiceResponse struct {
	ID               uint       `json:"id"`               // 請求書ID
	OrganizationID   uint       `json:"organizationId"`   // 請求元企業
	OrganizationName string     `json:"organizationName"` // 請求元企業名
	ClientID         uint       `json:"clientId"`         // 請求先取引先ID
	ClientName       string     `json:"clientName"`       // 請求先取引先名
	IssueDate        CustomDate `json:"issueDate"`        // 発行日
	Amount           int64      `json:"amount"`           // 請求金額
	Fee              int64      `json:"fee"`              // 手数料
	FeeRate          float64    `json:"feeRate"`          // 手数料率
	Tax              int64      `json:"tax"`              // 消費税
	TaxRate          float64    `json:"taxRate"`          // 消費税率
	TotalAmount      int64      `json:"totalAmount"`      // 合計金額
	DueDate          CustomDate `json:"dueDate"`          // 支払期日
	Status           string     `json:"status"`           // ステータス
}

type CustomDate struct {
	time.Time
}

func (d *CustomDate) UnmarshalJSON(b []byte) error {
	layout := `"2006-01-02"`
	parsedTime, err := time.Parse(layout, string(b))
	if err != nil {
		return err
	}
	d.Time = parsedTime
	return nil
}

func (d CustomDate) MarshalJSON() ([]byte, error) {
	const layout = "2006-01-02"
	return []byte(`"` + d.Time.Format(layout) + `"`), nil
}

func (h *InvoiceHandler) CreateInvoice(c echo.Context) error {
	var req CreateInvoiceRequest
	if err := c.Bind(&req); err != nil {
		log.Printf("Failed to bind request Error: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	invoice := application.CreateInvoiceDto{
		UserID:    req.UserID,
		ClientID:  req.ClientID,
		IssueDate: req.IssueDate.Time,
		Amount:    req.Amount,
		DueDate:   req.DueDate.Time,
	}

	createdInvoice, err := h.usecase.CreateInvoice(invoice)
	if err != nil {
		log.Printf("Failed to create invoice Error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not create invoice"})
	}

	response := CreateInvoiceResponse{
		ID:               createdInvoice.ID,
		OrganizationID:   createdInvoice.OrganizationID,
		OrganizationName: createdInvoice.OrganizationName,
		ClientID:         createdInvoice.ClientID,
		ClientName:       createdInvoice.ClientName,
		IssueDate:        CustomDate{createdInvoice.IssueDate},
		Amount:           createdInvoice.Amount,
		Fee:              createdInvoice.Fee,
		FeeRate:          createdInvoice.FeeRate,
		Tax:              createdInvoice.Tax,
		TaxRate:          createdInvoice.TaxRate,
		TotalAmount:      createdInvoice.TotalAmount,
		DueDate:          CustomDate{createdInvoice.DueDate},
		Status:           createdInvoice.Status,
	}

	return c.JSON(http.StatusOK, response)
}

type ListInvoiceRequest struct {
	StartDate int64  `query:"startDate"` // UnixTime
	EndDate   int64  `query:"endDate"`   // UnixTime
	TimeZone  string `query:"timeZone"`
}

type InvoiceItem struct {
	ID               uint       `json:"id"`               // 請求書ID
	OrganizationID   uint       `json:"organizationId"`   // 請求元企業
	OrganizationName string     `json:"organizationName"` // 請求元企業名
	ClientID         uint       `json:"clientId"`         // 請求先取引先ID
	ClientName       string     `json:"clientName"`       // 請求先取引先名
	IssueDate        CustomDate `json:"issueDate"`        // 発行日
	Amount           int64      `json:"amount"`           // 請求金額
	Fee              int64      `json:"fee"`              // 手数料
	FeeRate          float64    `json:"feeRate"`          // 手数料率
	Tax              int64      `json:"tax"`              // 消費税
	TaxRate          float64    `json:"taxRate"`          // 消費税率
	TotalAmount      int64      `json:"totalAmount"`      // 合計金額
	DueDate          CustomDate `json:"dueDate"`          // 支払期日
	Status           string     `json:"status"`           // ステータス
}

type ListInvoiceResponse struct {
	Invoices []InvoiceItem `json:"invoices" `
}

func (h *InvoiceHandler) ListInvoice(c echo.Context) error {
	var req ListInvoiceRequest
	if err := c.Bind(&req); err != nil {
		log.Printf("Failed to bind request Error: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	location, err := time.LoadLocation(req.TimeZone)
	if err != nil {
		log.Printf("Invalid timezone Error: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid timezone"})
	}

	// UNIXタイムを指定されたタイムゾーンに基づいて変換
	startDate := time.Unix(req.StartDate, 0).In(location)
	endDate := time.Unix(req.EndDate, 0).In(location)

	dto := application.ListInvoiceDto{
		StartDate: startDate,
		EndDate:   endDate,
	}

	invoices, err := h.usecase.ListInvoice(dto)
	if err != nil {
		log.Printf("Failed to list invoices Error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not list invoices"})
	}

	// DTOからレスポンスデータへの変換
	response := ListInvoiceResponse{
		Invoices: make([]InvoiceItem, len(invoices)),
	}

	for i, invoice := range invoices {
		response.Invoices[i] = InvoiceItem{
			ID:               invoice.ID,
			OrganizationID:   invoice.OrganizationID,
			OrganizationName: invoice.OrganizationName,
			ClientID:         invoice.ClientID,
			ClientName:       invoice.ClientName,
			IssueDate:        CustomDate{invoice.IssueDate},
			Amount:           invoice.Amount,
			Fee:              invoice.Fee,
			FeeRate:          invoice.FeeRate,
			Tax:              invoice.Tax,
			TaxRate:          invoice.TaxRate,
			TotalAmount:      invoice.TotalAmount,
			DueDate:          CustomDate{invoice.DueDate},
			Status:           invoice.Status,
		}
	}

	return c.JSON(http.StatusOK, response)
}
