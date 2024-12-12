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
	taxRateRepo      repository.TaxRate
}

func NewInvoiceUsecase(
	invoiceRepo repository.Invoice,
	clientRepo repository.Client,
	organizationRepo repository.Organization,
	taxRateRepo repository.TaxRate,
) *InvoiceUsecase {
	return &InvoiceUsecase{
		invoiceRepo:      invoiceRepo,
		clientRepo:       clientRepo,
		organizationRepo: organizationRepo,
		taxRateRepo:      taxRateRepo,
	}
}

type CreateInvoiceDto struct {
	UserID    uint
	ClientID  uint
	IssueDate time.Time
	Amount    int64
	DueDate   time.Time
}

type CreatedInvoiceDto struct {
	ID               uint
	OrganizationID   uint
	OrganizationName string
	ClientID         uint
	ClientName       string
	IssueDate        time.Time
	Amount           int64
	Fee              int64
	FeeRate          float64
	Tax              int64
	TaxRate          float64
	TotalAmount      int64
	DueDate          time.Time
	Status           string
}

// CreateInvoice 請求書を作成する.
// 現時点ではユースケース層に実装.
// ロジックを再利用したい場合や複雑になった場合はドメインサービスを作ることを検討する.
func (s *InvoiceUsecase) CreateInvoice(invoice CreateInvoiceDto) (*CreatedInvoiceDto, error) {
	// 会社を取得
	organization, err := s.organizationRepo.GetByID(invoice.UserID)
	if err != nil {
		return nil, err
	}

	// 取引先を取得
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

	// 消費税率を取得して金額を計算
	taxRate, err := s.taxRateRepo.GetRateByDate(invoice.IssueDate)
	if err != nil {
		return nil, err
	}
	newInvoice.Calculate(taxRate)

	// 請求書作成
	createdInvoice, err := s.invoiceRepo.Create(newInvoice)
	if err != nil {
		return nil, err
	}

	// dtoに変換
	dto, err := s.invoiceToDto(createdInvoice)
	if err != nil {
		return nil, err
	}

	return dto, nil
}

func (s *InvoiceUsecase) invoiceToDto(invoice *model.Invoice) (*CreatedInvoiceDto, error) {
	return &CreatedInvoiceDto{
		ID:               invoice.ID,
		OrganizationID:   invoice.Organization.ID,
		OrganizationName: invoice.Organization.Name,
		ClientID:         invoice.Client.ID,
		ClientName:       invoice.Client.Name,
		IssueDate:        invoice.IssueDate,
		Amount:           invoice.AmountAsInt(),
		Fee:              invoice.FeeAsInt(),
		FeeRate:          invoice.FeeRate,
		Tax:              invoice.TaxAsInt(),
		TaxRate:          invoice.TaxRate,
		TotalAmount:      invoice.TotalAmountAsInt(),
		DueDate:          invoice.DueDate,
		Status:           string(invoice.Status),
	}, nil
}

type ListInvoiceDto struct {
	StartDate time.Time
	EndDate   time.Time
}

func (s *InvoiceUsecase) ListInvoice(dto ListInvoiceDto) ([]*CreatedInvoiceDto, error) {
	// 指定された日付範囲内の請求書を取得
	invoices, err := s.invoiceRepo.FindByDueDateRange(dto.StartDate, dto.EndDate)
	if err != nil {
		return nil, err
	}

	// DTOリストに変換
	result := make([]*CreatedInvoiceDto, len(invoices))
	for i, invoice := range invoices {
		dto, err := s.invoiceToDto(invoice)
		if err != nil {
			return nil, err
		}
		result[i] = dto
	}

	return result, nil
}
