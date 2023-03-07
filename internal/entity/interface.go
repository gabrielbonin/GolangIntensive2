package entity

type OrderRepositoryInterface interface {
	Save(order *Order) (*Order, error)
	GetTotal() (int, error)
}
