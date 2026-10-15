package sale

import (
	"fmt"
	domainSale "inventario/internal/domain/sale"
)

type ListUseCase struct {
	saleRepo domainSale.Repository
}

func NewListUseCase(saleRepo domainSale.Repository) *ListUseCase {
	return &ListUseCase{saleRepo: saleRepo}
}

func (uc *ListUseCase) Execute() ([]domainSale.Sale, error) {
	sales, err := uc.saleRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("listar ventas: %w", err)
	}
	return sales, nil
}
