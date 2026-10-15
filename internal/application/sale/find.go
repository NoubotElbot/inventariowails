package sale

import (
	"fmt"
	domainSale "inventario/internal/domain/sale"
)

type FindUseCase struct {
	saleRepo domainSale.Repository
}

func NewFindUseCase(saleRepo domainSale.Repository) *FindUseCase {
	return &FindUseCase{saleRepo: saleRepo}
}

func (uc *FindUseCase) Execute(id int64) (*domainSale.Sale, error) {
	s, err := uc.saleRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("buscar venta: %w", err)
	}
	return s, nil
}
