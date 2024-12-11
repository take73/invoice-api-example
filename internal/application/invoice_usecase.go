package application

import (
	"time"

	"github.com/take73/invoice-api-example/internal/domain/model"
	"github.com/take73/invoice-api-example/internal/domain/repository"
)

type InvoiceUsecase struct {
	repo repository.InvoiceRepository
}

func NewInvoiceService(repo repository.InvoiceRepository) *InvoiceUsecase {
	return &InvoiceUsecase{repo: repo}
}

type CreateInvoiceDto struct {
	UserID    int
	IssueDate time.Time
	Amount    float64
	DueDate   time.Time
}

type CreatedInvoiceDto struct {
	ID               int       `json:"id"`               // 請求書ID
	OrganizationID   int       `json:"organizationId"`   // 請求元企業
	OrganizationName string    `json:"organizationName"` // 請求元企業名
	ClientID         int       `json:"clientId"`         // 請求先取引先ID
	ClientName       string    `json:"clientName"`       // 請求先取引先名
	IssueDate        time.Time `json:"issueDate"`        // 発行日
	Amount           float64   `json:"amount"`           // 請求金額
	Fee              float64   `json:"fee"`              // 手数料
	FeeRate          float64   `json:"feeRate"`          // 手数料率
	Tax              float64   `json:"tax"`              // 消費税
	TaxRate          float64   `json:"taxRate"`          // 消費税率
	TotalAmount      float64   `json:"totalAmount"`      // 合計金額
	DueDate          time.Time `json:"dueDate"`          // 支払期日
	Status           string    `json:"status"`           // ステータス
}

func (s *InvoiceUsecase) CreateInvoice(invoice CreateInvoiceDto) (*CreatedInvoiceDto, error) {
	in := model.Invoice{
		// Organization: invoice.,
		IssueDate: invoice.IssueDate,
		Amount:    invoice.Amount,
		DueDate:   invoice.DueDate,
	}

	out, err := s.repo.Create(in)
	if err != nil {
		return nil, err
	}

	dto := &CreatedInvoiceDto{
		ID:               out.ID,
		OrganizationID:   out.Organization.ID,
		OrganizationName: out.Organization.Name,
		ClientID:         out.Client.ID,
		ClientName:       out.Client.Name,
		IssueDate:        out.IssueDate,
		Amount:           out.Amount,
		Fee:              out.Fee,
		FeeRate:          out.FeeRate,
		Tax:              out.Tax,
		TaxRate:          out.TaxRate,
		TotalAmount:      out.TotalAmount,
		DueDate:          out.DueDate,
		Status:           string(out.Status),
	}
	return dto, nil
}
