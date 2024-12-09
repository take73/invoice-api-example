package application

type InvoiceUsecase struct {
    repo InvoiceRepository
}

type InvoiceRepository interface {
    CreateInvoice(data interface{}) error
}

func NewInvoiceService(repo InvoiceRepository) *InvoiceUsecase {
    return &InvoiceUsecase{repo: repo}
}

func (s *InvoiceUsecase) CreateInvoice(data interface{}) error {
    return nil
}