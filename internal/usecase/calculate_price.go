package usecase

import "github.com/gointensivo2/internal/entity"

type OrderInputDTO struct {
	ID    string
	Price float64
	Tax   float64
}

type OrderOutputDTO struct {
	ID         string
	Price      float64
	Tax        float64
	FinalPrice float64
}

type CalculateFinalPrice struct {
	OrderRepository entity.OrderRepositoryInterface
}

func (c *CalculateFinalPrice) Execute(input OrderInputDTO) (*OrderOutputDTO, error) {
	order, error := entity.NewOrder(input.ID, input.Price, input.Tax)

	if error != nil {
		return nil, error
	}

	error = order.CalculateFinalPrice()

	if error != nil {
		return nil, error
	}

	error = c.OrderRepository.Save(order)
	if error != nil {
		return nil, error
	}

	return &OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}, nil
}
