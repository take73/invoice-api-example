package http

import (
	"log"
	"net/http"
	"time"

	"github.com/take73/invoice-api-example/internal/application"

	"github.com/labstack/echo/v4"
)

const dateFormat = "2006-01-02"

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
	InvoiceItem
}

type CustomDate struct {
	time.Time
}

// / UnmarshalJSON JSONフィールドから日付をデコード
func (d *CustomDate) UnmarshalJSON(b []byte) error {
	str := string(b)
	// JSON の場合、クオートで囲まれているので削除
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}
	return d.unmarshalCommon(str)
}

// UnmarshalParam クエリパラメータから日付をデコード
func (d *CustomDate) UnmarshalParam(param string) error {
	return d.unmarshalCommon(param)
}

// unmarshalCommon 実際のパース処理を共通化
func (d *CustomDate) unmarshalCommon(dateStr string) error {
	if dateStr == "" {
		return nil // 空文字の場合はスキップ
	}
	parsedTime, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		return err
	}
	d.Time = parsedTime
	return nil
}

// MarshalJSON 日付をJSON形式でエンコード
func (d CustomDate) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil // ゼロ値の場合は null を返す
	}
	return []byte(`"` + d.Time.Format(dateFormat) + `"`), nil
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
		InvoiceItem: InvoiceItem{
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
		},
	}

	return c.JSON(http.StatusOK, response)
}

type ListInvoiceRequest struct {
	StartDate CustomDate `query:"startDate"`
	EndDate   CustomDate `query:"endDate"`
}

// 一旦postとgetで使いまわし
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
