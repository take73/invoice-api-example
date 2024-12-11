package application

import (
	"time"

	"github.com/take73/invoice-api-example/internal/domain/model"
	"github.com/take73/invoice-api-example/internal/domain/repository"
)

type InvoiceUsecase struct {
	invoiceRepo      repository.Invoice
	clientRepo       repository.Client
	organizationRepo repository.Organization
}

func NewInvoiceUsecase(
	invoiceRepo repository.Invoice,
	clientRepo repository.Client,
	organizationRepo repository.Organization,
) *InvoiceUsecase {
	return &InvoiceUsecase{invoiceRepo: invoiceRepo, clientRepo: clientRepo, organizationRepo: organizationRepo}
}

type CreateInvoiceDto struct {
	UserID    uint
	ClientID  uint
	IssueDate time.Time
	Amount    float64
	DueDate   time.Time
}

type CreatedInvoiceDto struct {
	ID               uint      `json:"id"`               // 請求書ID
	OrganizationID   uint      `json:"organizationId"`   // 請求元企業
	OrganizationName string    `json:"organizationName"` // 請求元企業名
	ClientID         uint      `json:"clientId"`         // 請求先取引先ID
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

// CreateInvoice 請求書を作成する.
// 現時点ではユースケース層に実装.
// ロジックを再利用したい場合や複雑になった場合はドメインサービスを作ることを検討する.
func (s *InvoiceUsecase) CreateInvoice(invoice CreateInvoiceDto) (*CreatedInvoiceDto, error) {
	organization, err := s.organizationRepo.GetByID(invoice.UserID)
	if err != nil {
		return nil, err
	}

	client, err := s.clientRepo.GetByID(invoice.ClientID)
	if err != nil {
		return nil, err
	}

	newInvoice, err := model.NewInvoice(
		organization,
		client,
		invoice.Amount,
		invoice.IssueDate,
		invoice.DueDate,
	)
	if err != nil {
		return nil, err
	}



	createdInvoice, err := s.invoiceRepo.Create(newInvoice)
	if err != nil {
		return nil, err
	}

	return &CreatedInvoiceDto{
		ID:               createdInvoice.ID,
		OrganizationID:   createdInvoice.Organization.ID,
		OrganizationName: createdInvoice.Organization.Name,
		ClientID:         createdInvoice.Client.ID,
		ClientName:       createdInvoice.Client.Name,
		IssueDate:        createdInvoice.IssueDate,
		Amount:           createdInvoice.Amount,
		Fee:              createdInvoice.Fee,
		FeeRate:          createdInvoice.FeeRate,
		Tax:              createdInvoice.Tax,
		TaxRate:          createdInvoice.TaxRate,
		TotalAmount:      createdInvoice.TotalAmount,
		DueDate:          createdInvoice.DueDate,
		Status:           string(createdInvoice.Status),
	}, nil
}
