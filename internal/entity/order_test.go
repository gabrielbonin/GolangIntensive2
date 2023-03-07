package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIfIdIsBlank(t *testing.T) {

	//err := order.Validate()
	//if err == nil {
	//t.Error("expected error")
	//}
	order := Order{}
	assert.Error(t, order.Validate(), "invalid id")

}

func TestIfPriceIsBlank(t *testing.T) {

	//err := order.Validate()
	//if err == nil {
	//t.Error("expected error")
	//}
	order := Order{ID: "123"}
	assert.Error(t, order.Validate(), "price is required")

}

func TestIfTaxIsBlank(t *testing.T) {

	//err := order.Validate()
	//if err == nil {
	//t.Error("expected error")
	//}
	order := Order{ID: "123", Price: 10.0}
	assert.Error(t, order.Validate(), "tax is required")

}

func TestIfAllValidParams(t *testing.T) {

	//err := order.Validate()
	//if err == nil {
	//t.Error("expected error")
	//}
	order := Order{ID: "123", Price: 10.0, Tax: 1.0}
	assert.NoError(t, order.Validate())
	assert.Equal(t, 10.0, order.Price)
	assert.Equal(t, 1.0, order.Tax)
	order.CalculateFinalPrice()
	assert.Equal(t, 11.0, order.FinalPrice)
}
