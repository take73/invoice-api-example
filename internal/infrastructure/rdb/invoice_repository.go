package rdb

type InvoiceRepository struct{}

func NewInvoiceRepository() *InvoiceRepository {
	return &InvoiceRepository{}
}

func (r *InvoiceRepository) CreateInvoice(data interface{}) error {
	// データ保存ロジックを追加する予定
	return nil
}
