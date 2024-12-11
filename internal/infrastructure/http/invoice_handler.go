package http

import (
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
	UserID    int        `json:"userId"`
	ClientID  int        `json:"clientId"`
	IssueDate CustomDate `json:"issueDate"`
	Amount    float64    `json:"amount"`
	DueDate   CustomDate `json:"dueDate"`
}

type CreateInvoiceResponse struct {
	ID               int        `json:"id"`               // 請求書ID
	OrganizationID   int        `json:"organizationId"`   // 請求元企業
	OrganizationName string     `json:"organizationName"` // 請求元企業名
	ClientID         int        `json:"clientId"`         // 請求先取引先ID
	ClientName       string     `json:"clientName"`       // 請求先取引先名
	IssueDate        CustomDate `json:"issueDate"`        // 発行日
	Amount           float64    `json:"amount"`           // 請求金額
	Fee              float64    `json:"fee"`              // 手数料
	FeeRate          float64    `json:"feeRate"`          // 手数料率
	Tax              float64    `json:"tax"`              // 消費税
	TaxRate          float64    `json:"taxRate"`          // 消費税率
	TotalAmount      float64    `json:"totalAmount"`      // 合計金額
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
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	invoice := application.CreateInvoiceDto{
		UserID:    req.UserID,
		IssueDate: req.IssueDate.Time,
		Amount:    req.Amount,
		DueDate:   req.DueDate.Time,
	}

	createdInvoice, err := h.usecase.CreateInvoice(invoice)
	if err != nil {
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
