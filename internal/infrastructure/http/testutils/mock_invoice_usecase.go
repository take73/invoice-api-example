package testutils

import (
	"github.com/stretchr/testify/mock"
	"github.com/take73/invoice-api-example/internal/application"
)

type MockInvoiceUsecase struct {
	mock.Mock
}

func (m *MockInvoiceUsecase) CreateInvoice(dto application.CreateInvoiceDto) (*application.CreatedInvoiceDto, error) {
	args := m.Called(dto)
	if args.Get(0) != nil {
		return args.Get(0).(*application.CreatedInvoiceDto), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockInvoiceUsecase) ListInvoice(dto application.ListInvoiceDto) ([]*application.CreatedInvoiceDto, error) {
	args := m.Called(dto)
	if args.Get(0) != nil {
		return args.Get(0).([]*application.CreatedInvoiceDto), args.Error(1)
	}
	return nil, args.Error(1)
}
