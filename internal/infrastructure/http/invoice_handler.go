package http

import (
	"log"
	"net/http"

	"github.com/take73/invoice-api-example/internal/application"
	"github.com/take73/invoice-api-example/internal/shared/types"

	"github.com/labstack/echo/v4"
)

type InvoiceHandler struct {
	usecase application.InvoiceUsecase
}

func NewInvoiceHandler(usecase application.InvoiceUsecase) *InvoiceHandler {
	return &InvoiceHandler{usecase: usecase}
}

type CreateInvoiceRequest struct {
	UserID    uint             `json:"userId" validate:"required,gt=0"`   // 必須, 0より大きい
	ClientID  uint             `json:"clientId" validate:"required,gt=0"` // 必須, 0より大きい
	IssueDate types.CustomDate `json:"issueDate" validate:"required"`     // 必須, 有効な日付
	Amount    int64            `json:"amount"`                            // 0やマイナスを許容しないなら必須としてもよさそう
	DueDate   types.CustomDate `json:"dueDate" validate:"required"`       // 必須, 有効な日付
}

type CreateInvoiceResponse struct {
	InvoiceItem
}

func (h *InvoiceHandler) CreateInvoice(c echo.Context) error {
	var req CreateInvoiceRequest
	if err := c.Bind(&req); err != nil {
		log.Printf("Failed to bind request Error: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	// Validate the request payload
	if err := c.Validate(&req); err != nil {
		log.Printf("Validation failed Error: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "validation failed"})
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
		InvoiceItem: InvoiceItem{
			ID:               createdInvoice.ID,
			OrganizationID:   createdInvoice.OrganizationID,
			OrganizationName: createdInvoice.OrganizationName,
			ClientID:         createdInvoice.ClientID,
			ClientName:       createdInvoice.ClientName,
			IssueDate:        types.CustomDate{Time: createdInvoice.IssueDate},
			Amount:           createdInvoice.Amount,
			Fee:              createdInvoice.Fee,
			FeeRate:          createdInvoice.FeeRate,
			Tax:              createdInvoice.Tax,
			TaxRate:          createdInvoice.TaxRate,
			TotalAmount:      createdInvoice.TotalAmount,
			DueDate:          types.CustomDate{Time: createdInvoice.DueDate},
			Status:           createdInvoice.Status,
		},
	}

	return c.JSON(http.StatusOK, response)
}

type ListInvoiceRequest struct {
	StartDate types.CustomDate `query:"startDate"`
	EndDate   types.CustomDate `query:"endDate"`
}

// 一旦postとgetで使いまわし
type InvoiceItem struct {
	ID               uint             `json:"id"`               // 請求書ID
	OrganizationID   uint             `json:"organizationId"`   // 請求元企業
	OrganizationName string           `json:"organizationName"` // 請求元企業名
	ClientID         uint             `json:"clientId"`         // 請求先取引先ID
	ClientName       string           `json:"clientName"`       // 請求先取引先名
	IssueDate        types.CustomDate `json:"issueDate"`        // 発行日
	Amount           int64            `json:"amount"`           // 請求金額
	Fee              int64            `json:"fee"`              // 手数料
	FeeRate          float64          `json:"feeRate"`          // 手数料率
	Tax              int64            `json:"tax"`              // 消費税
	TaxRate          float64          `json:"taxRate"`          // 消費税率
	TotalAmount      int64            `json:"totalAmount"`      // 合計金額
	DueDate          types.CustomDate `json:"dueDate"`          // 支払期日
	Status           string           `json:"status"`           // ステータス
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

	dto := application.ListInvoiceDto{
		StartDate: req.StartDate.Time,
		EndDate:   req.EndDate.Time,
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
			IssueDate:        types.CustomDate{Time: invoice.IssueDate},
			Amount:           invoice.Amount,
			Fee:              invoice.Fee,
			FeeRate:          invoice.FeeRate,
			Tax:              invoice.Tax,
			TaxRate:          invoice.TaxRate,
			TotalAmount:      invoice.TotalAmount,
			DueDate:          types.CustomDate{Time: invoice.DueDate},
			Status:           invoice.Status,
		}
	}

	return c.JSON(http.StatusOK, response)
}
